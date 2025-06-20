package engine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/ext/resource/loader"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/loopfz/gadgeto/tonic"
	"k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/resource"
)

func newEngineHandler(cl cluster.Cluster, config APIConfiguration) (gin.HandlerFunc, error) {
	openAPI, err := cl.OpenAPIClient("1.32")
	if err != nil {
		return nil, err
	}

	policyLoader, err := loader.New(openAPI)
	if err != nil {
		return nil, err
	}

	tcm := patch.NewTypeConverterManager(nil, openAPI)
	go tcm.Run(context.Background())

	return tonic.Handler(func(ctx *gin.Context, in *EngineRequest) (*EngineResponse, error) {
		params, err := in.LoadParameters()
		if err != nil {
			return nil, fmt.Errorf("unable to load params: %w", err)
		}
		params.ImageData = in.ImageData
		policies, vaps, vapbs, vpols, ivpols, dpols, gpols, mpols, err := in.LoadPolicies(policyLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load policies: %w", err)
		}
		vapbsWindow, err := in.LoadVAPBindings(policyLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load policies: %w", err)
		}
		vapbs = append(vapbs, vapbsWindow...)

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

		clusterResources = append(oldResources, clusterResources...)
		clusterObjects := resource.AppendNamespaces(resources, clusterResources)

		dClient, err := cl.DClient(resource.ToObjects(resources), clusterObjects...)
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

		processor, err := engine.NewProcessor(params, cl, config, dClient, cmResolver, exceptionSelector, tcm)
		if err != nil {
			return nil, err
		}
		results, err := processor.Run(ctx, policies, vaps, vapbs, vpols, ivpols, dpols, gpols, mpols, resources, oldResources)
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

				return fmt.Errorf("variable %s is not defined in the context", variable.Name)
			}
		}
	}

	return nil
}
