package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/store"
	"github.com/kyverno/kyverno/pkg/config"
	"github.com/kyverno/kyverno/pkg/engine"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type engineRequest struct {
	Policy    string `json:"policy"`
	Resources string `json:"resources"`
	Context   string `json:"context"`
}

type engineResponse struct {
	Validation []engineapi.EngineResponse
}

func (r engineRequest) loadResources() ([]unstructured.Unstructured, error) {
	if documents, err := yamlutils.SplitDocuments([]byte(r.Resources)); err != nil {
		return nil, err
	} else {
		var resources []unstructured.Unstructured
		for _, document := range documents {
			var resource unstructured.Unstructured
			if resourceJson, err := yaml.YAMLToJSON(document); err != nil {
				return nil, err
			} else if err := resource.UnmarshalJSON(resourceJson); err != nil {
				return nil, err
			}
			resources = append(resources, resource)
		}
		return resources, nil
	}
}
func (r engineRequest) process(ctx context.Context) (*engineResponse, error) {
	if policies, err := yamlutils.GetPolicy([]byte(r.Policy)); err != nil {
		return nil, err
	} else if resources, err := r.loadResources(); err != nil {
		return nil, err
	} else {
		var response engineResponse
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
				// TODO: set data in engine context
				policyContext := engine.NewPolicyContextWithJsonContext(kyvernov1.Create, engineContext).
					WithPolicy(policy).
					WithNewResource(resource)
				// WithNamespaceLabels(namespaceLabels).
				// WithAdmissionInfo(c.UserInfo).
				// WithResourceKind(gvk, subresource)
				response.Validation = append(response.Validation, eng.Validate(ctx, policyContext))
			}
		}
		return &response, nil
	}
}

func run(c *gin.Context) {
	// TODO: error handling
	var request engineRequest
	if err := c.BindJSON(&request); err != nil {
		return
	} else if response, err := request.process(c.Request.Context()); err != nil {
		return
	} else {
		c.IndentedJSON(http.StatusNotModified, response)
	}
}

func main() {
	var host = flag.String("host", "localhost", "server host")
	var port = flag.Int("port", 8080, "server port")
	var frontendPath = flag.String("frontend-path", "../frontend/dist", "frontend folder")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	router.POST("/engine", run)
	router.StaticFS("/", http.Dir(*frontendPath))
	address := fmt.Sprintf("%v:%v", *host, *port)
	if err := router.Run(address); err != nil {
		panic(err)
	}
}
