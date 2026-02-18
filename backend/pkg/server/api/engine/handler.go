package engine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/loopfz/gadgeto/tonic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/engine/cel/processor"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/resource"
)

func newEngineHandler(cl cluster.Cluster, config APIConfiguration) (gin.HandlerFunc, error) {
	return tonic.Handler(func(ctx *gin.Context, in *EngineRequest) (*EngineResponse, error) {
		params, err := in.LoadParameters()
		if err != nil {
			return nil, fmt.Errorf("unable to load params: %w", err)
		}
		params.ImageData = in.ImageData

		openAPI, err := cl.OpenAPIClient(params.Kubernetes.Version)
		if err != nil {
			return nil, err
		}

		policyLoader, err := loader.New(openAPI)
		if err != nil {
			return nil, err
		}

		client, err := in.OpenAPIClient(cl, params.Kubernetes.Version, config)
		if err != nil {
			return nil, fmt.Errorf("failed to load OpenAPI client: %w", err)
		}

		c, cancel := context.WithCancel(context.Background())
		defer cancel()

		tcm := patch.NewTypeConverterManager(nil, client)
		go tcm.Run(c)

		resourceLoader, err := in.ResourceLoader(client)
		if err != nil {
			return nil, err
		}

		k8s, jsonPolicies, authzPolicies, err := in.LoadPolicies(policyLoader)
		if err != nil {
			return nil, fmt.Errorf("unable to load policies: %w", err)
		}

		var results *models.Results
		var resources []unstructured.Unstructured

		if k8s.Length() > 0 {
			vapbsWindow, err := in.LoadVAPBindings(policyLoader)
			if err != nil {
				return nil, fmt.Errorf("unable to load policies: %w", err)
			}
			k8s.ValidatingAdmissionPolicyBindings = append(k8s.ValidatingAdmissionPolicyBindings, vapbsWindow...)

			resources, err = resource.LoadResources(resourceLoader, []byte(in.Resources))
			if err != nil {
				return nil, fmt.Errorf("unable to load resources: %w", err)
			}
			oldResources, err := resource.LoadResources(resourceLoader, []byte(in.OldResources))
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
			crds, err := in.LoadCRDs()
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
				if err := validateParams(params, cmResolver, k8s.Policies); err != nil {
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
			resources, err = resource.LoadJSON(in.Resources)
			if err != nil {
				return nil, fmt.Errorf("unable to load resources: %w", err)
			}

			dClient, err := CreateDClient(cl, in, resourceLoader, nil)
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
			resources, err := resource.LoadEnvyRequests(in.Resources)
			if err != nil {
				return nil, fmt.Errorf("unable to load resources: %w", err)
			}

			dClient, err := CreateDClient(cl, in, resourceLoader, nil)
			if err != nil {
				return nil, err
			}

			processor := processor.NewAuthz(dClient)
			results, err = processor.RunEnvoy(ctx, authzPolicies, resources)
			if err != nil {
				return nil, err
			}

			return &EngineResponse{
				Results: results,
			}, nil
		}

		if len(authzPolicies.HTTPPolicies) > 0 {
			resources, err := resource.LoadHTTPRequests(in.Resources)
			if err != nil {
				return nil, fmt.Errorf("unable to load resources: %w", err)
			}

			dClient, err := CreateDClient(cl, in, resourceLoader, nil)
			if err != nil {
				return nil, err
			}

			processor := processor.NewAuthz(dClient)
			results, err = processor.RunHTTP(ctx, authzPolicies, resources)
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
	}, http.StatusOK), nil
}

func CreateDClient(cl cluster.Cluster, in *EngineRequest, loader loader.Loader, resources []unstructured.Unstructured) (dclient.Interface, error) {
	clusterResources, err := in.LoadClusterResources(loader)
	if err != nil {
		return nil, err
	}

	clusterObjects := resource.AppendNamespaces(resources, clusterResources)

	return cl.DClient(clusterObjects)
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
