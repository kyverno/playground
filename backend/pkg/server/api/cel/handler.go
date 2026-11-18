package cel

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	gctxstore "github.com/kyverno/kyverno/pkg/globalcontext/store"
	"github.com/loopfz/gadgeto/tonic"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine/cel/expr"
	"github.com/kyverno/playground/backend/pkg/engine/models"
)

// EvaluateRequest is the payload for evaluating a standalone CEL expression
// against the same variables and libraries a ValidatingPolicy would have.
type EvaluateRequest struct {
	// Expression is the CEL expression to evaluate.
	Expression string `json:"expression"`
	// Resource is an optional YAML/JSON manifest bound to the `object` variable.
	Resource string `json:"resource,omitempty"`
	// OldResource is an optional YAML/JSON manifest bound to the `oldObject` variable.
	OldResource string `json:"oldResource,omitempty"`
	// NamespaceResource is an optional YAML/JSON Namespace manifest bound to the `namespaceObject` variable.
	NamespaceResource string `json:"namespaceResource,omitempty"`
	// Context carries the admission context (operation, user info, ...) used to build the `request` variable.
	Context *models.Context `json:"context,omitempty"`
}

// EvaluateResponse is the outcome of evaluating a standalone CEL expression.
// A failure to compile or evaluate the expression is reported in Error
// rather than as an HTTP error, so it can be rendered next to the
// expression like any other result.
type EvaluateResponse struct {
	Value any    `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
	Error string `json:"error,omitempty"`
}

func newHandler(cl cluster.Cluster) (gin.HandlerFunc, error) {
	return tonic.Handler(func(ctx *gin.Context, in *EvaluateRequest) (*EvaluateResponse, error) {
		resource, err := parseResource(in.Resource)
		if err != nil {
			return nil, fmt.Errorf("failed to parse resource: %w", err)
		}
		oldResource, err := parseResource(in.OldResource)
		if err != nil {
			return nil, fmt.Errorf("failed to parse oldResource: %w", err)
		}
		namespaceResource, err := parseResource(in.NamespaceResource)
		if err != nil {
			return nil, fmt.Errorf("failed to parse namespaceResource: %w", err)
		}

		admCtx := models.Context{Operation: kyvernov1.Create}
		if in.Context != nil {
			admCtx = *in.Context
		}

		dClient, err := cl.DClient(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}
		contextProvider, err := libs.NewContextProvider(dClient, nil, gctxstore.New(0), cl.RESTMapper(nil), false)
		if err != nil {
			return nil, fmt.Errorf("failed to create CEL library context: %w", err)
		}

		result, err := expr.Evaluate(ctx, contextProvider, in.Expression, expr.Input{
			Object:          resource,
			OldObject:       oldResource,
			Request:         buildAdmissionRequest(resource, oldResource, admCtx),
			NamespaceObject: namespaceResource,
			Namespace:       resourceNamespace(resource, oldResource),
		})
		if err != nil {
			return nil, err
		}

		return &EvaluateResponse{Value: result.Value, Type: result.Type, Error: result.Error}, nil
	}, http.StatusOK), nil
}

func parseResource(content string) (*unstructured.Unstructured, error) {
	if strings.TrimSpace(content) == "" {
		return nil, nil
	}
	var obj map[string]any
	if err := yaml.Unmarshal([]byte(content), &obj); err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: obj}, nil
}

func resourceNamespace(resources ...*unstructured.Unstructured) string {
	for _, r := range resources {
		if r != nil && r.GetNamespace() != "" {
			return r.GetNamespace()
		}
	}
	return ""
}

func buildAdmissionRequest(resource, oldResource *unstructured.Unstructured, admCtx models.Context) *admissionv1.AdmissionRequest {
	primary := resource
	if primary == nil {
		primary = oldResource
	}
	if primary == nil {
		primary = &unstructured.Unstructured{}
	}
	gvk := primary.GroupVersionKind()
	gvr := gvk.GroupVersion().WithResource(strings.ToLower(gvk.Kind) + "s")

	var newBytes, oldBytes []byte
	if resource != nil {
		newBytes, _ = resource.MarshalJSON()
	}
	if oldResource != nil {
		oldBytes, _ = oldResource.MarshalJSON()
	}

	dryRun := admCtx.DryRun
	return &admissionv1.AdmissionRequest{
		UID:       "abc-123",
		Kind:      metav1.GroupVersionKind(gvk),
		Resource:  metav1.GroupVersionResource(gvr),
		Name:      primary.GetName(),
		Namespace: primary.GetNamespace(),
		Operation: admissionv1.Operation(admCtx.Operation),
		UserInfo: authenticationv1.UserInfo{
			UID:      "user-123",
			Username: admCtx.Username,
			Groups:   admCtx.Groups,
		},
		Object:    runtime.RawExtension{Raw: newBytes},
		OldObject: runtime.RawExtension{Raw: oldBytes},
		DryRun:    &dryRun,
	}
}
