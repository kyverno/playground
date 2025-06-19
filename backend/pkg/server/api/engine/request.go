package engine

import (
	"errors"
	"os"
	"testing/fstest"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2 "github.com/kyverno/kyverno/api/kyverno/v2"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/ext/resource/loader"
	v1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/exception"
	"github.com/kyverno/playground/backend/pkg/policy"
	"github.com/kyverno/playground/backend/pkg/resource"
	"github.com/kyverno/playground/backend/pkg/vapbinding"
)

type EngineRequest struct {
	Policies                  string                      `json:"policies"`
	Resources                 string                      `json:"resources"`
	OldResources              string                      `json:"oldResources"`
	ClusterResources          string                      `json:"clusterResources"`
	Context                   string                      `json:"context"`
	Config                    string                      `json:"config"`
	CustomResourceDefinitions string                      `json:"customResourceDefinitions"`
	PolicyExceptions          string                      `json:"policyExceptions"`
	VAPBindings               string                      `json:"vapBindings"`
	ImageData                 map[string]models.ImageData `json:"imageData"`
}

func (r *EngineRequest) LoadParameters() (*models.Parameters, error) {
	var params models.Parameters
	if err := yaml.Unmarshal([]byte(r.Context), &params); err != nil {
		return nil, err
	}
	return &params, nil
}

func (r *EngineRequest) LoadPolicies(policyLoader loader.Loader) ([]kyvernov1.PolicyInterface, []v1.ValidatingAdmissionPolicy, []v1.ValidatingAdmissionPolicyBinding, []v1alpha1.ValidatingPolicy, []v1alpha1.ImageValidatingPolicy, []v1alpha1.DeletingPolicy, []v1alpha1.GeneratingPolicy, []v1alpha1.MutatingPolicy, error) {
	return policy.Load(policyLoader, []byte(r.Policies))
}

func (r *EngineRequest) LoadResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	data, err := resource.LoadJSON(r.Resources)
	if errors.Is(err, resource.ErrNoJSON) {
		return resource.LoadResources(resourceLoader, []byte(r.Resources))
	}
	return data, err
}

func (r *EngineRequest) LoadClusterResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	return resource.LoadResources(resourceLoader, []byte(r.ClusterResources))
}

func (r *EngineRequest) LoadOldResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	data, err := resource.LoadJSON(r.OldResources)
	if errors.Is(err, resource.ErrNoJSON) {
		return resource.LoadResources(resourceLoader, []byte(r.OldResources))
	}

	return data, err
}

func (r *EngineRequest) LoadPolicyExceptions() ([]*kyvernov2.PolicyException, error) {
	return exception.Load([]byte(r.PolicyExceptions))
}

func (r *EngineRequest) LoadVAPBindings(policyLoader loader.Loader) ([]v1.ValidatingAdmissionPolicyBinding, error) {
	return vapbinding.Load(policyLoader, []byte(r.VAPBindings))
}

func (r *EngineRequest) LoadConfig(resourceLoader loader.Loader) (*corev1.ConfigMap, error) {
	if len(r.Config) == 0 {
		return nil, nil
	}
	return resource.Load[corev1.ConfigMap](resourceLoader, []byte(r.Config))
}

func (r *EngineRequest) ResourceLoader(cluster cluster.Cluster, kubeVersion string, config APIConfiguration) (loader.Loader, error) {
	var clients []openapi.Client
	if cluster != nil && !cluster.IsFake() {
		dclient, err := cluster.DClient(nil)
		if err != nil {
			return nil, err
		}
		clients = append(clients, dclient.GetKubeClient().Discovery().OpenAPIV3())
	} else {
		client, err := cluster.OpenAPIClient(kubeVersion)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	schemas, err := data.Schemas()
	if err != nil {
		return nil, err
	}

	clients = append(clients, openapiclient.NewLocalSchemaFiles(schemas))
	if len(r.CustomResourceDefinitions) != 0 {
		mapFs := fstest.MapFS{
			"crds.yaml": &fstest.MapFile{
				Data: []byte(r.CustomResourceDefinitions),
			},
		}
		clients = append(clients, openapiclient.NewLocalCRDFiles(mapFs))
	}
	for _, crd := range config.LocalCrds {
		clients = append(clients, openapiclient.NewLocalCRDFiles(os.DirFS(crd)))
	}
	for _, crd := range config.BuiltInCrds {
		fs, err := data.BuiltInCrds(crd)
		if err != nil {
			return nil, err
		}

		clients = append(clients, openapiclient.NewLocalCRDFiles(fs))
	}
	return loader.New(openapiclient.NewComposite(clients...))
}
