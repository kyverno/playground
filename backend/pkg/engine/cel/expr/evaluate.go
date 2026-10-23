package expr

import (
	"context"
	gojson "encoding/json"
	"fmt"
	"reflect"

	"github.com/kyverno/kyverno/pkg/cel/compiler"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	sdkutils "github.com/kyverno/sdk/extensions/cel/utils"
	"google.golang.org/protobuf/types/known/structpb"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/cel/lazy"
)

// Result is the outcome of evaluating a standalone CEL expression. Exactly
// one of Value or Error is set.
type Result struct {
	// Value is the JSON-serializable result of the expression.
	Value any `json:"value,omitempty"`
	// Type is the CEL type name of the result (e.g. "bool", "string", "list").
	Type string `json:"type,omitempty"`
	// Error is set when the expression failed to compile or evaluate.
	Error string `json:"error,omitempty"`
}

// Input holds the variables a standalone CEL expression can reference,
// matching what a ValidatingPolicy sees: object, oldObject, request and
// namespaceObject.
type Input struct {
	Object          *unstructured.Unstructured
	OldObject       *unstructured.Unstructured
	Request         *admissionv1.AdmissionRequest
	NamespaceObject *unstructured.Unstructured
	Namespace       string
}

// Evaluate compiles and evaluates a standalone CEL expression against the
// same variables and CEL libraries a Kyverno ValidatingPolicy has access to.
// Compile and evaluation errors are returned inside Result rather than as a
// Go error, so callers can render them next to the expression like any other
// evaluation outcome; a non-nil error return indicates the environment
// itself could not be constructed.
func Evaluate(ctx context.Context, libsctx libs.Context, expression string, in Input) (*Result, error) {
	env, err := newEnv(libsctx, in.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to build CEL environment: %w", err)
	}

	ast, issues := env.Compile(expression)
	if err := issues.Err(); err != nil {
		return &Result{Error: err.Error()}, nil
	}

	program, err := env.Program(ast)
	if err != nil {
		return &Result{Error: err.Error()}, nil
	}

	requestVal, err := sdkutils.ConvertObjectToUnstructured(in.Request)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request variable: %w", err)
	}

	data := map[string]any{
		compiler.ObjectKey:          objectOrNil(in.Object),
		compiler.OldObjectKey:       objectOrNil(in.OldObject),
		compiler.RequestKey:         requestVal.Object,
		compiler.NamespaceObjectKey: objectOrNil(in.NamespaceObject),
		compiler.VariablesKey:       lazy.NewMapValue(compiler.VariablesType),
		compiler.ExceptionsKey:      libs.Exception{AllowedImages: []string{}, AllowedValues: []string{}},
	}

	out, _, err := program.ContextEval(ctx, data)
	if err != nil {
		return &Result{Error: err.Error()}, nil
	}

	native, err := out.ConvertToNative(reflect.TypeOf(&structpb.Value{}))
	if err != nil {
		return &Result{Error: fmt.Sprintf("result could not be converted to JSON: %s", err)}, nil
	}

	raw, err := gojson.Marshal(native)
	if err != nil {
		return &Result{Error: fmt.Sprintf("result could not be marshalled to JSON: %s", err)}, nil
	}

	var value any
	if err := gojson.Unmarshal(raw, &value); err != nil {
		return &Result{Error: fmt.Sprintf("result could not be marshalled to JSON: %s", err)}, nil
	}

	return &Result{
		Value: value,
		Type:  out.Type().TypeName(),
	}, nil
}

func objectOrNil(u *unstructured.Unstructured) any {
	if u == nil {
		return nil
	}
	return u.Object
}
