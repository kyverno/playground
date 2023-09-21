package engine

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/loopfz/gadgeto/tonic"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine"
)

func newEngineHandler(cl cluster.Cluster, config APIConfiguration) (gin.HandlerFunc, error) {
	policyClient := openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas")
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
		exceptions, err := in.LoadPolicyExceptions(policyLoader)
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
		processor, err := engine.NewProcessor(params, cl, config, dClient, cmResolver, exceptionSelector)
		if err != nil {
			return nil, err
		}
		results, err := processor.Run(ctx, policies, vaps, resources, oldResources)
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
		return "1.27", nil
	}
	version, err := semver.NewVersion(kubeVersion)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(version.Major(), ".", version.Minor()), nil
}
