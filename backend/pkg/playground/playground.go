package playground

import (
	"context"
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/engine/cel/processor"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/policy"
	"github.com/kyverno/playground/backend/pkg/resource"
)

func Run(ctx context.Context, cl cluster.Cluster, request *EngineRequest, config crd.APIConfiguration) (*EngineResponse, error) {
	params, err := request.LoadParameters()
	if err != nil {
		return nil, err
	}

	openAPI, err := cl.OpenAPIClient(params.Kubernetes.Version)
	if err != nil {
		return nil, err
	}

	policyLoader, err := loader.New(openAPI)
	if err != nil {
		return nil, err
	}

	client, err := crd.OpenAPIClient(cl, params.Kubernetes.Version, "", config)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI client: %w", err)
	}

	resourceLoader, err := request.ResourceLoader(client)
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	tcm := patch.NewTypeConverterManager(nil, client)
	go tcm.Run(c)

	k8s, jsonPolicies, authzPolicies, err := policy.Load(policyLoader, []byte(request.Policies))
	if err != nil {
		return nil, err
	}

	var results *models.Results
	var resources []unstructured.Unstructured

	if k8s.Length() > 0 {
		resources, err = resource.LoadResources(resourceLoader, []byte(request.Resources))
		if err != nil {
			return nil, fmt.Errorf("unable to load resources: %w", err)
		}
		oldResources, err := resource.LoadResources(resourceLoader, []byte(request.OldResources))
		if err != nil {
			return nil, fmt.Errorf("unable to load old resources: %w", err)
		}
		clusterResources, err := resource.LoadResources(resourceLoader, []byte(request.ClusterResources))
		if err != nil {
			return nil, err
		}
		config, err := request.LoadConfig(resourceLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load config resources: %w", err)
		}
		exceptions, celExceptions, err := request.LoadPolicyExceptions()
		if err != nil {
			return nil, fmt.Errorf("unable to load policy exceptions: %w", err)
		}

		k8s.PolicyExceptions = append(k8s.PolicyExceptions, celExceptions...)
		jsonPolicies.PolicyExceptions = append(jsonPolicies.PolicyExceptions, celExceptions...)
		authzPolicies.PolicyExceptions = append(authzPolicies.PolicyExceptions, celExceptions...)

		crds, err := request.LoadCRDs()
		if err != nil {
			return nil, fmt.Errorf("unable to load resources: %w", err)
		}

		clusterResources = append(oldResources, clusterResources...)
		clusterObjects := resource.AppendNamespaces(resources, clusterResources)

		dClient, err := cl.DClient(append(resource.ToObjects(resources), clusterObjects...))
		if err != nil {
			return nil, err
		}
		cmResolver, err := cluster.NewConfigMapResolver(dClient)
		if err != nil {
			return nil, err
		}

		var exceptionSelector engineapi.PolicyExceptionSelector
		if params.Flags.Exceptions.Enabled {
			exceptionSelector = cl.PolicyExceptionSelector(params.Flags.Exceptions.Namespace, exceptions...)
		}

		if cl.IsFake() {
			if err := ValidateParams(params, cmResolver, k8s.Policies); err != nil {
				return nil, err
			}
		}

		processor, err := engine.NewProcessor(params, cl, config, dClient, cmResolver, exceptionSelector, tcm, cl.RESTMapper(crds))
		if err != nil {
			return nil, err
		}
		results, err = processor.Run(ctx, k8s, resources, oldResources)
		if err != nil {
			return nil, err
		}

		return &EngineResponse{
			Resources: resources,
			Results:   results,
		}, nil
	}

	if jsonPolicies.Length() > 0 {
		resources, err = resource.LoadJSON(request.Resources)
		if err != nil {
			return nil, fmt.Errorf("unable to load resources: %w", err)
		}

		dClient, err := CreateDClient(cl, request, resourceLoader, nil)
		if err != nil {
			return nil, err
		}

		processor := processor.NewJSON(dClient, tcm)
		results, err = processor.Run(ctx, jsonPolicies, resources)
		if err != nil {
			return nil, err
		}

		return &EngineResponse{
			Resources: resources,
			Results:   results,
		}, nil
	}

	if len(authzPolicies.EnvoyPolicies) > 0 {
		resources, err := resource.LoadEnvyRequests(request.Resources)
		if err != nil {
			return nil, fmt.Errorf("unable to load resources: %w", err)
		}

		dClient, err := CreateDClient(cl, request, resourceLoader, nil)
		if err != nil {
			return nil, err
		}

		processor := processor.NewAuthz(dClient)
		results, err = processor.RunEnvoy(context.Background(), authzPolicies, resources)
		if err != nil {
			return nil, err
		}

		return &EngineResponse{
			Results: results,
		}, nil
	}

	if len(authzPolicies.HTTPPolicies) > 0 {
		resources, err := resource.LoadHTTPRequests(request.Resources)
		if err != nil {
			return nil, fmt.Errorf("unable to load resources: %w", err)
		}

		dClient, err := CreateDClient(cl, request, resourceLoader, nil)
		if err != nil {
			return nil, err
		}

		processor := processor.NewAuthz(dClient)
		results, err = processor.RunHTTP(context.Background(), authzPolicies, resources)
		if err != nil {
			return nil, err
		}

		return &EngineResponse{
			Results: results,
		}, nil
	}

	return &EngineResponse{
		Resources: nil,
		Results:   nil,
	}, nil
}

func ValidateParams(params *models.Parameters, cmResolver engineapi.ConfigmapResolver, policies []kyvernov1.PolicyInterface) error {
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

func CreateDClient(cl cluster.Cluster, in *EngineRequest, loader loader.Loader, resources []unstructured.Unstructured) (dclient.Interface, error) {
	clusterResources, err := in.LoadClusterResources(loader)
	if err != nil {
		return nil, err
	}

	clusterObjects := resource.AppendNamespaces(resources, clusterResources)

	return cl.DClient(clusterObjects)
}
