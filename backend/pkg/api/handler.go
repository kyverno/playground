package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/context/resolvers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/resource/loader"
	"github.com/kyverno/playground/backend/pkg/utils"
)

func NewEngineHandler(cluster cluster.Cluster) (gin.HandlerFunc, error) {
	var kubeClient kubernetes.Interface
	var dClient dclient.Interface
	var cmResolver engineapi.ConfigmapResolver
	var resourceLoader loader.Loader
	if cluster != nil {
		kubeClient = cluster.KubeClient()
		dClient = cluster.DClient()
		loader, err := loader.NewRemote(cluster)
		if err != nil {
			return nil, err
		}
		resourceLoader = loader
	}
	if kubeClient != nil {
		resolver, err := resolvers.NewClientBasedResolver(kubeClient)
		if err != nil {
			return nil, err
		}
		cmResolver = resolver
	}
	return func(c *gin.Context) {
		var request EngineRequest
		err := c.ShouldBind(&request)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		params, err := engine.ParseParameters(request.Context)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid context string")
			return
		}

		if resourceLoader == nil {
			loader, err := loader.NewLocal(params.Kubernetes.Version)
			if err != nil {
				c.String(http.StatusInternalServerError, "failed to initialize loader")
				return
			}
			resourceLoader = loader
		}

		resources, err := loader.LoadResources(resourceLoader, []byte(request.Resources))
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "failed to parse resources")
			return
		}

		policies, err := utils.LoadPolicies(resourceLoader, []byte(request.Policies))
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to parse policies")
			return
		}

		var config *corev1.ConfigMap
		if len(request.Config) != 0 {
			conf, err := loader.Load[corev1.ConfigMap](resourceLoader, []byte(request.Config))
			if err != nil {
				c.String(http.StatusInternalServerError, "failed to parse kyverno configmap")
				return
			}
			config = conf
		}

		processor, err := engine.NewProcessor(params, config, dClient, cmResolver)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to initialize processor")
			return
		}

		results, err := processor.Run(c, policies, resources)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		response := &EngineResponse{
			Policies:          policies,
			Resources:         resources,
			Mutation:          results.Mutation,
			ImageVerification: results.ImageVerification,
			Validation:        results.Validation,
			Generation:        results.Generation,
		}

		c.IndentedJSON(http.StatusOK, response)
	}, nil
}
