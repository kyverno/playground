package engine

import (
	"os"
	"testing/fstest"

	kyvernov2 "github.com/kyverno/kyverno/api/kyverno/v2"
	"github.com/kyverno/kyverno/ext/resource/loader"
	yamlutils "github.com/kyverno/kyverno/ext/yaml"
	v1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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

func (r *EngineRequest) LoadPolicies(policyLoader loader.Loader) (policy.K8sPolicies, policy.JSONPolicies, policy.AuthzPolicies, error) {
	return policy.Load(policyLoader, []byte(r.Policies))
}

func (r *EngineRequest) LoadCRDs() ([]*apiextensionsv1.CustomResourceDefinition, error) {
	scheme := runtime.NewScheme()
	if err := apiextensionsv1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	documents, err := yamlutils.SplitDocuments([]byte(r.CustomResourceDefinitions))
	if err != nil {
		return nil, err
	}

	crds := []*apiextensionsv1.CustomResourceDefinition{}
	for _, document := range documents {
		crd := &apiextensionsv1.CustomResourceDefinition{}
		if err := yaml.Unmarshal(document, crd); err != nil {
			return nil, err
		}
		crds = append(crds, crd)
	}
	return crds, nil
}

func (r *EngineRequest) LoadClusterResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	return resource.LoadResources(resourceLoader, []byte(r.ClusterResources))
}

func (r *EngineRequest) LoadOldResources(resourceLoader loader.Loader) ([]unstructured.Unstructured, error) {
	return resource.LoadResources(resourceLoader, []byte(r.OldResources))
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

func (r *EngineRequest) OpenAPIClient(cluster cluster.Cluster, kubeVersion string, config APIConfiguration) (openapi.Client, error) {
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

	return openapiclient.NewComposite(clients...), nil
}

func (r *EngineRequest) ResourceLoader(client openapi.Client) (loader.Loader, error) {
	return loader.New(client)
}
