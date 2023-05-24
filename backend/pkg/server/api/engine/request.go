package engine

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/resource/loader"
	"github.com/kyverno/playground/backend/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type EngineRequest struct {
	Policies  string `json:"policies"`
	Resources string `json:"resources"`
	Context   string `json:"context"`
	Config    string `json:"config"`
}

func (r *EngineRequest) ParseContext() (*engine.Parameters, error) {
	var params engine.Parameters
	if err := yaml.Unmarshal([]byte(r.Context), &params); err != nil {
		return nil, err
	}
	return &params, nil
}

func (r *EngineRequest) LoadPolicies(policyLoader loader.Loader) ([]kyvernov1.PolicyInterface, error) {
	return utils.LoadPolicies(policyLoader, []byte(r.Policies))
}

func (r *EngineRequest) LoadResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	return loader.LoadResources(resourceLoader, []byte(r.Resources))
}

func (r *EngineRequest) LoadConfig(resourceLoader loader.Loader) (*corev1.ConfigMap, error) {
	if len(r.Config) == 0 {
		return nil, nil
	}
	return loader.Load[corev1.ConfigMap](resourceLoader, []byte(r.Config))
}
