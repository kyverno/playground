package engine

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/resource/loader"
	"github.com/kyverno/playground/backend/pkg/utils"
)

type EngineRequest struct {
	Policies                  string `json:"policies"`
	Resources                 string `json:"resources"`
	OldResources              string `json:"oldResources"`
	ClusterResources          string `json:"clusterResources"`
	Context                   string `json:"context"`
	Config                    string `json:"config"`
	CustomResourceDefinitions string `json:"customResourceDefinitions"`
	PolicyExceptions          string `json:"policyExceptions"`
}

func (r *EngineRequest) LoadParameters() (*engine.Parameters, error) {
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

func (r *EngineRequest) LoadClusterResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	return loader.LoadResources(resourceLoader, []byte(r.ClusterResources))
}

func (r *EngineRequest) LoadOldResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	return loader.LoadResources(resourceLoader, []byte(r.OldResources))
}

func (r *EngineRequest) LoadPolicyExceptions(resourceLoader loader.Loader) ([]*kyvernov2alpha1.PolicyException, error) {
	return utils.LoadPolicyExceptions(resourceLoader, []byte(r.PolicyExceptions))
}

func (r *EngineRequest) LoadConfig(resourceLoader loader.Loader) (*corev1.ConfigMap, error) {
	if len(r.Config) == 0 {
		return nil, nil
	}
	return loader.Load[corev1.ConfigMap](resourceLoader, []byte(r.Config))
}

func (r *EngineRequest) ResourceLoader(kubeVersion string, config APIConfiguration) (loader.Loader, error) {
	kubeVersion, err := parseKubeVersion(kubeVersion)
	if err != nil {
		return nil, err
	}
	clients := []openapi.Client{
		openapiclient.NewHardcodedBuiltins(kubeVersion),
	}
	if len(r.CustomResourceDefinitions) != 0 {
		clients = append(clients, NewInMemory([]byte(r.CustomResourceDefinitions)))
	}
	for _, crd := range config.LocalCrds {
		clients = append(clients, openapiclient.NewLocalCRDFiles(nil, crd))
	}
	for _, crd := range config.BuiltInCrds {
		fs, path := data.BuiltInCrds(crd)
		clients = append(clients, openapiclient.NewLocalCRDFiles(fs, path))
	}
	return loader.New(openapiclient.NewComposite(clients...))
}
