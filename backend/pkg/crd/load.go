package crd

import (
	"os"
	"testing/fstest"

	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/cluster"
)

type APIConfiguration struct {
	BuiltInCrds []string
	LocalCrds   []string
}

func OpenAPIClient(cluster cluster.Cluster, kubeVersion, customResourceDefinitions string, config APIConfiguration) (openapi.Client, error) {
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
	if len(customResourceDefinitions) != 0 {
		mapFs := fstest.MapFS{
			"crds.yaml": &fstest.MapFile{
				Data: []byte(customResourceDefinitions),
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
