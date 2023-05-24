package engine

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/context/resolvers"
	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/resource/loader"
	"github.com/loopfz/gadgeto/tonic"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

func newEngineHandler(cluster cluster.Cluster) (gin.HandlerFunc, error) {
	policyClient := openapiclient.NewLocalFiles(data.Schemas(), "schemas")
	policyLoader, err := loader.New(policyClient)
	if err != nil {
		return nil, err
	}
	return tonic.Handler(func(c *gin.Context, in *EngineRequest) (*EngineResponse, error) {
		params, err := in.ParseContext()
		if err != nil {
			return nil, err
		}
		policies, err := in.LoadPolicies(policyLoader)
		if err != nil {
			return nil, err
		}
		resourceLoader, err := resourceLoader(params.Kubernetes.Version)
		if err != nil {
			return nil, err
		}
		resources, err := in.LoadResources(resourceLoader)
		if err != nil {
			return nil, err
		}
		config, err := in.LoadConfig(resourceLoader)
		if err != nil {
			return nil, err
		}
		processor, err := getProcessor(params, config, cluster)
		if err != nil {
			return nil, err
		}
		results, err := processor.Run(c, policies, resources)
		if err != nil {
			return nil, err
		}
		return &EngineResponse{
			Policies:          policies,
			Resources:         resources,
			Mutation:          results.Mutation,
			ImageVerification: results.ImageVerification,
			Validation:        results.Validation,
			Generation:        results.Generation,
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

func resourceLoader(kubeVersion string) (loader.Loader, error) {
	kubeVersion, err := parseKubeVersion(kubeVersion)
	if err != nil {
		return nil, err
	}
	return loader.New(openapiclient.NewHardcodedBuiltins(kubeVersion))
}

func getProcessor(params *engine.Parameters, config *corev1.ConfigMap, cluster cluster.Cluster) (*engine.Processor, error) {
	var dClient dclient.Interface
	var cmResolver engineapi.ConfigmapResolver
	var exceptionSelector engineapi.PolicyExceptionSelector
	if cluster != nil {
		dClient = cluster.DClient()
		exceptionSelector = engine.NewPolicyExceptionSelector(cluster)
		if kubeClient := cluster.KubeClient(); kubeClient != nil {
			resolver, err := resolvers.NewClientBasedResolver(kubeClient)
			if err != nil {
				return nil, err
			}
			cmResolver = resolver
		}
	}
	return engine.NewProcessor(params, config, dClient, cmResolver, exceptionSelector)
}
