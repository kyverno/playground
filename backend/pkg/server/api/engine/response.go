package engine

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
)

type EngineResponse struct {
	Resources []unstructured.Unstructured `json:"resources"`
	*models.Results
}
