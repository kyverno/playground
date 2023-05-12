package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov1beta1 "github.com/kyverno/kyverno/api/kyverno/v1beta1"
	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/store"
	"github.com/kyverno/kyverno/pkg/config"
	"github.com/kyverno/kyverno/pkg/engine"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	"github.com/kyverno/kyverno/pkg/engine/policycontext"
	"github.com/kyverno/kyverno/pkg/registryclient"
	jsonutils "github.com/kyverno/kyverno/pkg/utils/json"
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	"github.com/kyverno/playground/backend/data"
	_ "go.etcd.io/etcd/client/pkg/v3/logutil"
	authenticationv1 "k8s.io/api/authentication/v1"
	_ "k8s.io/apiextensions-apiserver/pkg/apiserver/conversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/kubectl-validate/pkg/validatorfactory"
	"sigs.k8s.io/yaml"
)

type apiContext struct {
	Username        string                       `json:"username"`
	Groups          []string                     `json:"groups"`
	Roles           []string                     `json:"roles"`
	ClusterRoles    []string                     `json:"clusterRoles"`
	Operation       kyvernov1.AdmissionOperation `json:"operation"`
	NamespaceLabels map[string]string            `json:"namespaceLabels"`
}

type apiRequest struct {
	Policies  string `json:"policies"`
	Resources string `json:"resources"`
	Context   string `json:"context"`
}

type apiResponse struct {
	Policies          []kyvernov1.PolicyInterface `json:"policies"`
	Resources         []unstructured.Unstructured `json:"resources"`
	Mutation          []EngineResponse            `json:"mutation"`
	ImageVerification []EngineResponse            `json:"imageVerification"`
	Validation        []EngineResponse            `json:"validation"`
	Generation        []EngineResponse            `json:"generation"`
}

type EngineResponse struct {
	// OriginalResource is the original resource as YAML string
	OriginalResource string `json:"originalResource"`
	// Resource is the original resource
	Resource unstructured.Unstructured `json:"resource"`
	// Policy is the original policy
	Policy kyvernov1.PolicyInterface `json:"policy"`
	// namespaceLabels given by policy context
	NamespaceLabels map[string]string `json:"namespaceLabels"`
	// PatchedResource is the resource patched with the engine action changes
	PatchedResource string `json:"patchedResource"`
	// PolicyResponse contains the engine policy response
	PolicyResponse PolicyResponse `json:"policyResponse"`
}

type PolicyResponse struct {
	// Rules contains policy rules responses
	Rules []RuleResponse `json:"rules"`
}

type RuleResponse struct {
	// name is the rule name specified in policy
	Name string `json:"name"`
	// ruleType is the rule type (Mutation,Generation,Validation) for Kyverno Policy
	RuleType engineapi.RuleType `json:"ruleType"`
	// message is the message response from the rule application
	Message string `json:"message"`
	// status rule status
	Status engineapi.RuleStatus `json:"status"`
	// patches are JSON patches, for mutation rules
	Patches []string `json:"patches"`
	// generatedResource is the generated by the generate rules of a policy
	GeneratedResource unstructured.Unstructured `json:"generatedResource"`
	// patchedTarget is the patched resource for mutate.targets
	PatchedTarget *unstructured.Unstructured `json:"patchedTarget"`
	// patchedTargetParentResourceGVR is the GVR of the parent resource of the PatchedTarget. This is only populated when PatchedTarget is a subresource.
	PatchedTargetParentResourceGVR metav1.GroupVersionResource `json:"patchedTargetParentResourceGVR"`
	// patchedTargetSubresourceName is the name of the subresource which is patched, empty if the resource patched is not a subresource.
	PatchedTargetSubresourceName string `json:"patchedTargetSubresourceName"`
	// podSecurityChecks contains pod security checks (only if this is a pod security rule)
	PodSecurityChecks *engineapi.PodSecurityChecks `json:"podSecurityChecks"`
	// exception is the exception applied (if any)
	Exception *kyvernov2alpha1.PolicyException `json:"exception"`
}

func ConvertRuleResponse(in engineapi.RuleResponse) RuleResponse {
	out := RuleResponse{
		Name:              in.Name(),
		RuleType:          in.RuleType(),
		Message:           in.Message(),
		Status:            in.Status(),
		GeneratedResource: in.GeneratedResource(),
		// PatchedTarget *unstructured.Unstructured
		// // patchedTargetParentResourceGVR is the GVR of the parent resource of the PatchedTarget. This is only populated when PatchedTarget is a subresource.
		// PatchedTargetParentResourceGVR metav1.GroupVersionResource
		// // patchedTargetSubresourceName is the name of the subresource which is patched, empty if the resource patched is not a subresource.
		// PatchedTargetSubresourceName string
		PodSecurityChecks: in.PodSecurityChecks(),
		Exception:         in.Exception(),
	}
	for _, patch := range in.Patches() {
		out.Patches = append(out.Patches, string(patch))
	}
	return out
}

func ConvertEngineResponse(in engineapi.EngineResponse) EngineResponse {
	patchedResource, _ := yaml.Marshal(in.PatchedResource.Object)
	resource, _ := yaml.Marshal(in.Resource.Object)
	out := EngineResponse{
		OriginalResource: string(resource),
		Resource:         in.Resource,
		Policy:           in.Policy(),
		NamespaceLabels:  in.NamespaceLabels(),
		PatchedResource:  string(patchedResource),
	}
	for _, ruleresponse := range in.PolicyResponse.Rules {
		out.PolicyResponse.Rules = append(out.PolicyResponse.Rules, ConvertRuleResponse(ruleresponse))
	}
	return out
}

func loadUnstructured(document []byte) (unstructured.Unstructured, error) {
	const mediaType = runtime.ContentTypeYAML
	var result unstructured.Unstructured
	var metadata metav1.TypeMeta
	if err := yaml.Unmarshal(document, &metadata); err != nil {
		return result, err
	}
	gvk := metadata.GetObjectKind().GroupVersionKind()
	if gvk.Empty() {
		return result, fmt.Errorf("GVK cannot be empty")
	}
	if factory, err := validatorfactory.New(
		openapiclient.NewComposite(
			openapiclient.NewLocalFiles(data.Schemas(), "schemas"),
			openapiclient.NewHardcodedBuiltins("1.27"),
		),
	); err != nil {
		return result, err
	} else if validator, err := factory.ValidatorsForGVK(gvk); err != nil {
		return result, err
	} else if decoder, err := validator.Decoder(gvk); err != nil {
		return result, err
	} else if info, ok := runtime.SerializerInfoForMediaType(decoder.SupportedMediaTypes(), mediaType); !ok {
		return result, fmt.Errorf("unsupported media type %q", mediaType)
	} else if _, _, err := decoder.DecoderToVersion(info.StrictSerializer, gvk.GroupVersion()).Decode(document, &gvk, &result); err != nil {
		return result, err
	} else {
		return result, nil
	}
}

func fromUnstructured[T any](untyped unstructured.Unstructured) (T, error) {
	var result T
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(untyped.UnstructuredContent(), &result); err != nil {
		return result, err
	}
	return result, nil
}

func (r apiRequest) loadPolicies() ([]kyvernov1.PolicyInterface, error) {
	loadPolicy := func(untyped unstructured.Unstructured) (kyvernov1.PolicyInterface, error) {
		kind := untyped.GetKind()
		if kind == "Policy" {
			if policy, err := fromUnstructured[kyvernov1.Policy](untyped); err != nil {
				return nil, err
			} else {
				return &policy, nil
			}
		} else if kind == "ClusterPolicy" {
			if policy, err := fromUnstructured[kyvernov1.ClusterPolicy](untyped); err != nil {
				return nil, err
			} else {
				return &policy, nil
			}
		} else {
			return nil, fmt.Errorf("invalid kind: %s", kind)
		}
	}
	if documents, err := yamlutils.SplitDocuments([]byte(r.Policies)); err != nil {
		return nil, err
	} else {
		var policies []kyvernov1.PolicyInterface
		for _, document := range documents {
			if untyped, err := loadUnstructured(document); err != nil {
				return nil, err
			} else if policy, err := loadPolicy(untyped); err != nil {
				return nil, err
			} else {
				policies = append(policies, policy)
			}
		}
		return policies, nil
	}
}

func (r apiRequest) loadResources() ([]unstructured.Unstructured, error) {
	if documents, err := yamlutils.SplitDocuments([]byte(r.Resources)); err != nil {
		return nil, err
	} else {
		var resources []unstructured.Unstructured
		for _, document := range documents {
			var resource unstructured.Unstructured
			if resourceJSON, err := yaml.YAMLToJSON(document); err != nil {
				return nil, err
			} else if err := resource.UnmarshalJSON(resourceJSON); err != nil {
				continue
			}
			resources = append(resources, resource)
		}
		return resources, nil
	}
}

func (r apiRequest) loadContext() (apiContext, error) {
	ctx := apiContext{}
	err := yaml.Unmarshal([]byte(r.Context), &ctx)

	return ctx, err
}

func (r apiRequest) process(ctx context.Context) (*apiResponse, error) {
	if policies, err := r.loadPolicies(); err != nil {
		return nil, err
	} else if resources, err := r.loadResources(); err != nil {
		return nil, err
	} else if requestContext, err := r.loadContext(); err != nil {
		return nil, err
	} else {
		apiResponse := apiResponse{
			Resources: resources,
			Policies:  policies,
		}
		cfg := config.NewDefaultConfiguration(false)
		jp := jmespath.New(cfg)
		rclient, err := registryclient.New(registryclient.WithLocalKeychain())
		if err != nil {
			return nil, err
		}
		engine := engine.NewEngine(
			cfg,
			config.NewDefaultMetricsConfiguration(),
			jp,
			nil,
			rclient,
			store.ContextLoaderFactory(nil),
			nil,
		)
		admissionInfo := kyvernov1beta1.RequestInfo{
			AdmissionUserInfo: authenticationv1.UserInfo{
				Username: requestContext.Username,
				Groups:   requestContext.Groups,
			},
			Roles:        requestContext.Roles,
			ClusterRoles: requestContext.ClusterRoles,
		}
		operation := requestContext.Operation
		if operation == "" {
			operation = kyvernov1.Create
		}
		for _, resource := range resources {
			resource := resource
			getContext := func(policy kyvernov1.PolicyInterface) (engineapi.PolicyContext, error) {
				policyContext, err := policycontext.NewPolicyContext(
					jp,
					resource,
					operation,
					&admissionInfo,
					cfg,
				)
				if err != nil {
					return nil, err
				}
				policyContext = policyContext.
					WithPolicy(policy).
					WithNamespaceLabels(requestContext.NamespaceLabels)
				// WithResourceKind(gvk, subresource)
				return policyContext, nil
			}
			// mutate
			for _, policy := range policies {
				if policyContext, err := getContext(policy); err != nil {
					return nil, err
				} else {
					response := engine.Mutate(ctx, policyContext)
					resource = response.PatchedResource
					apiResponse.Mutation = append(apiResponse.Mutation, ConvertEngineResponse(response))
				}
			}
			// verify images
			for _, policy := range policies {
				if policyContext, err := getContext(policy); err != nil {
					return nil, err
				} else {
					response, verifiedImageData := engine.VerifyAndPatchImages(ctx, policyContext)
					// TODO: we apply patches manually because the engine doesn't
					patches := response.GetPatches()
					if !verifiedImageData.IsEmpty() {
						annotationPatches, err := verifiedImageData.Patches(len(resource.GetAnnotations()) != 0, logr.Discard())
						if err != nil {
							return nil, err
						} else {
							// add annotation patches first
							patches = append(annotationPatches, patches...)
						}
					}
					if len(patches) != 0 {
						patch := jsonutils.JoinPatches(patches...)
						decoded, err := jsonpatch.DecodePatch(patch)
						if err != nil {
							return nil, err
						}
						options := &jsonpatch.ApplyOptions{SupportNegativeIndices: true, AllowMissingPathOnRemove: true, EnsurePathExistsOnAdd: true}
						resourceBytes, err := resource.MarshalJSON()
						if err != nil {
							return nil, err
						}
						patchedResourceBytes, err := decoded.ApplyWithOptions(resourceBytes, options)
						if err != nil {
							return nil, err
						}
						if err := response.PatchedResource.UnmarshalJSON(patchedResourceBytes); err != nil {
							return nil, err
						}
					}
					resource = response.PatchedResource
					apiResponse.ImageVerification = append(apiResponse.ImageVerification, ConvertEngineResponse(response))
				}
			}
			// validate
			for _, policy := range policies {
				if policyContext, err := getContext(policy); err != nil {
					return nil, err
				} else {
					response := engine.Validate(ctx, policyContext)
					apiResponse.Validation = append(apiResponse.Validation, ConvertEngineResponse(response))
				}
			}
			// generate
			for _, policy := range policies {
				if policyContext, err := getContext(policy); err != nil {
					return nil, err
				} else {
					response := engine.Generate(ctx, policyContext)
					apiResponse.Generation = append(apiResponse.Generation, ConvertEngineResponse(response))
				}
			}
		}
		return &apiResponse, nil
	}
}

func serveApi(c *gin.Context) {
	// TODO: error handling
	var request apiRequest
	if err := c.BindJSON(&request); err != nil {
		return
	} else if response, err := request.process(c.Request.Context()); err != nil {
		return
	} else {
		c.IndentedJSON(http.StatusOK, response)
	}
}

func run(host string, port int) {
	fs, err := fs.Sub(data.StaticFiles(), "dist")
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	router.POST("/engine", serveApi)
	router.StaticFS("/", http.FS(fs))
	address := fmt.Sprintf("%v:%v", host, port)
	if err := router.Run(address); err != nil {
		panic(err)
	}
}

func main() {
	var host = flag.String("host", "localhost", "server host")
	var port = flag.Int("port", 8080, "server port")
	var mode = flag.String("mode", gin.ReleaseMode, "gin run mode")
	flag.Parse()
	gin.SetMode(*mode)
	run(*host, *port)
}
