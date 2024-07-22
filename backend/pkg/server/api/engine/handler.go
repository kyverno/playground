package engine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-gonic/gin"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/ext/resource/loader"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/loopfz/gadgeto/tonic"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/engine/models"
)

func newEngineHandler(cl cluster.Cluster, config APIConfiguration) (gin.HandlerFunc, error) {
	policyClient := openapiclient.NewComposite(
		openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas"),
		openapiclient.NewGitHubBuiltins("1.28"),
	)
	policyLoader, err := loader.New(policyClient)
	if err != nil {
		return nil, err
	}
	return tonic.Handler(func(ctx *gin.Context, in *EngineRequest) (*EngineResponse, error) {
		params, err := in.LoadParameters()
		if err != nil {
			return nil, fmt.Errorf("unable to load params: %w", err)
		}
		params.ImageData = in.ImageData
		policies, vaps, err := in.LoadPolicies(policyLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load policies: %w", err)
		}
		vapbs, err := in.LoadVAPBindings(policyLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load policies: %w", err)
		}
		resourceLoader, err := in.ResourceLoader(cl, params.Kubernetes.Version, config)
		if err != nil {
			return nil, err
		}
		resources, err := in.LoadResources(resourceLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load resources: %w", err)
		}
		oldResources, err := in.LoadOldResources(resourceLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load old resources: %w", err)
		}
		clusterResources, err := in.LoadClusterResources(resourceLoader)
		if err != nil {
			return nil, err
		}
		config, err := in.LoadConfig(resourceLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load config resources: %w", err)
		}
		exceptions, err := in.LoadPolicyExceptions()
		if err != nil {
			return nil, fmt.Errorf("unable to load policy exceptions: %w", err)
		}
		dClient, err := cl.DClient(clusterResources...)
		if err != nil {
			return nil, err
		}
		cmResolver, err := cluster.NewConfigMapResolver(dClient)
		if err != nil {
			return nil, err
		}
		// TODO: move in engine ?
		var exceptionSelector engineapi.PolicyExceptionSelector
		if params.Flags.Exceptions.Enabled {
			exceptionSelector = cl.PolicyExceptionSelector(params.Flags.Exceptions.Namespace, exceptions...)
		}

		if cl.IsFake() {
			if err := validateParams(params, cmResolver, policies); err != nil {
				fmt.Println(err)
				return nil, err
			}
		}

		processor, err := engine.NewProcessor(params, cl, config, dClient, cmResolver, exceptionSelector)
		if err != nil {
			return nil, err
		}
		results, err := processor.Run(ctx, policies, vaps, vapbs, resources, oldResources)
		if err != nil {
			return nil, err
		}
		return &EngineResponse{
			Policies:  policies,
			Resources: resources,
			Results:   results,
		}, nil
	}, http.StatusOK), nil
}

func parseKubeVersion(kubeVersion string) (string, error) {
	if kubeVersion == "" {
		return "1.28", nil
	}
	version, err := semver.NewVersion(kubeVersion)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(version.Major(), ".", version.Minor()), nil
}

func validateParams(params *models.Parameters, cmResolver engineapi.ConfigmapResolver, policies []kyvernov1.PolicyInterface) error {
	if params == nil {
		return nil
	}

	for _, policy := range policies {
		for _, rule := range policy.GetSpec().Rules {
			for _, variable := range rule.Context {
				if variable.APICall == nil && variable.ConfigMap == nil {
					continue
				}
				if _, ok := params.Variables[variable.Name]; ok {
					continue
				}
				if variable.ConfigMap != nil {
					_, err := cmResolver.Get(context.Background(), variable.ConfigMap.Namespace, variable.ConfigMap.Name)
					if err == nil {
						continue
					}
				}

				return fmt.Errorf("Variable %s is not defined in the context", variable.Name)
			}
		}
	}

	return nil
}
