package expr

import (
	"context"
	"testing"

	"github.com/kyverno/kyverno/pkg/cel/libs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestEvaluate_Literal(t *testing.T) {
	result, err := Evaluate(context.Background(), libs.NewFakeContextProvider(), "1 + 1", Input{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Error != "" {
		t.Fatalf("unexpected evaluation error: %s", result.Error)
	}
	if result.Value != float64(2) {
		t.Fatalf("expected 2, got %v (%T)", result.Value, result.Value)
	}
	if result.Type != "int" {
		t.Fatalf("expected type int, got %s", result.Type)
	}
}

func TestEvaluate_ObjectField(t *testing.T) {
	object := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"name": "my-pod",
		},
	}}

	result, err := Evaluate(context.Background(), libs.NewFakeContextProvider(), "object.metadata.name", Input{Object: object})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Error != "" {
		t.Fatalf("unexpected evaluation error: %s", result.Error)
	}
	if result.Value != "my-pod" {
		t.Fatalf("expected my-pod, got %v", result.Value)
	}
}

func TestEvaluate_CompileError(t *testing.T) {
	result, err := Evaluate(context.Background(), libs.NewFakeContextProvider(), "object.metadata.", Input{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Error == "" {
		t.Fatal("expected a compile error to be reported")
	}
}

func TestEvaluate_KyvernoLibrary(t *testing.T) {
	// exercises one of the Kyverno CEL libraries (math) to confirm the
	// environment carries the same libraries a ValidatingPolicy would have.
	result, err := Evaluate(context.Background(), libs.NewFakeContextProvider(), `math.round(3.14159, 2)`, Input{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Error != "" {
		t.Fatalf("unexpected evaluation error: %s", result.Error)
	}
	if result.Value != 3.14 {
		t.Fatalf("expected 3.14, got %v", result.Value)
	}
}
