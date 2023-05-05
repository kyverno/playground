package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov1beta1 "github.com/kyverno/kyverno/api/kyverno/v1beta1"
	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/store"
	"github.com/kyverno/kyverno/pkg/config"
	"github.com/kyverno/kyverno/pkg/engine"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

//go:embed dist
var staticFiles embed.FS

type apiContext struct {
	Username        string                       `json:"policy"`
	Groups          []string                     `json:"groups"`
	Roles           []string                     `json:"roles"`
	ClusterRoles    []string                     `json:"clusterRoles"`
	Operation       kyvernov1.AdmissionOperation `json:"operation"`
	NamespaceLabels map[string]string            `json:"namespaceLabels"`
}

type apiRequest struct {
	Policy    string     `json:"policy"`
	Resources string     `json:"resources"`
	Context   apiContext `json:"context"`
}

type apiResponse struct {
	Policies   []kyvernov1.PolicyInterface
	Resources  []unstructured.Unstructured
	Validation []EngineResponse
}

type EngineResponse struct {
	// Resource is the original resource
	Resource unstructured.Unstructured `json:"resource"`
	// Policy is the original policy
	Policy kyvernov1.PolicyInterface `json:"policy"`
	// namespaceLabels given by policy context
	NamespaceLabels map[string]string `json:"namespaceLabels"`
	// PatchedResource is the resource patched with the engine action changes
	PatchedResource unstructured.Unstructured `json:"patchedResource"`
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
		Name:     in.Name(),
		RuleType: in.RuleType(),
		Message:  in.Message(),
		Status:   in.Status(),
		// Patches []string
		GeneratedResource: in.GeneratedResource(),
		// PatchedTarget *unstructured.Unstructured
		// // patchedTargetParentResourceGVR is the GVR of the parent resource of the PatchedTarget. This is only populated when PatchedTarget is a subresource.
		// PatchedTargetParentResourceGVR metav1.GroupVersionResource
		// // patchedTargetSubresourceName is the name of the subresource which is patched, empty if the resource patched is not a subresource.
		// PatchedTargetSubresourceName string
		PodSecurityChecks: in.PodSecurityChecks(),
		Exception:         in.Exception(),
	}
	return out
}

func ConvertEngineResponse(in engineapi.EngineResponse) EngineResponse {
	out := EngineResponse{
		Resource:        in.Resource,
		Policy:          in.Policy,
		NamespaceLabels: in.NamespaceLabels(),
		PatchedResource: in.PatchedResource,
	}
	for _, ruleresponse := range in.PolicyResponse.Rules {
		out.PolicyResponse.Rules = append(out.PolicyResponse.Rules, ConvertRuleResponse(ruleresponse))
	}
	return out
}

func (r apiRequest) loadResources() ([]unstructured.Unstructured, error) {
	if documents, err := yamlutils.SplitDocuments([]byte(r.Resources)); err != nil {
		return nil, err
	} else {
		var resources []unstructured.Unstructured
		for _, document := range documents {
			var resource unstructured.Unstructured
			if resourceJson, err := yaml.YAMLToJSON(document); err != nil {
				return nil, err
			} else if err := resource.UnmarshalJSON(resourceJson); err != nil {
				continue
			}
			resources = append(resources, resource)
		}
		return resources, nil
	}
}
func (r apiRequest) process(ctx context.Context) (*apiResponse, error) {
	if policies, err := yamlutils.GetPolicy([]byte(r.Policy)); err != nil {
		return nil, err
	} else if resources, err := r.loadResources(); err != nil {
		return nil, err
	} else {
		response := apiResponse{
			Resources: resources,
			Policies:  policies,
		}
		cfg := config.NewDefaultConfiguration(false)
		jp := jmespath.New(cfg)
		eng := engine.NewEngine(
			cfg,
			config.NewDefaultMetricsConfiguration(),
			jp,
			nil,
			nil,
			store.ContextLoaderFactory(nil),
			nil,
		)
		for _, resource := range resources {
			for _, policy := range policies {
				engineContext := enginecontext.NewContext(jp)
				operation := r.Context.Operation
				if operation == "" {
					operation = kyvernov1.Create
				}
				// TODO: set data in engine context
				policyContext := engine.NewPolicyContextWithJsonContext(operation, engineContext).
					WithPolicy(policy).
					WithNewResource(resource).
					WithNamespaceLabels(r.Context.NamespaceLabels).
					WithAdmissionInfo(kyvernov1beta1.RequestInfo{
						AdmissionUserInfo: authenticationv1.UserInfo{
							Username: r.Context.Username,
							Groups:   r.Context.Groups,
						},
						Roles:        r.Context.Roles,
						ClusterRoles: r.Context.ClusterRoles,
					})
				// WithResourceKind(gvk, subresource)
				response.Validation = append(response.Validation, ConvertEngineResponse(eng.Validate(ctx, policyContext)))
			}
		}
		return &response, nil
	}
}

func run(c *gin.Context) {
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

func main() {
	var host = flag.String("host", "localhost", "server host")
	var port = flag.Int("port", 8080, "server port")
	fs, err := fs.Sub(staticFiles, "dist")
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
	router.POST("/engine", run)
	router.StaticFS("/", http.FS(fs))
	address := fmt.Sprintf("%v:%v", *host, *port)
	if err := router.Run(address); err != nil {
		panic(err)
	}
}
