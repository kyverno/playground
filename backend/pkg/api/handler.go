package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	corev1 "k8s.io/api/core/v1"

	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/resource/loader"
	"github.com/kyverno/playground/backend/pkg/utils"
)

func NewEngineHandler(dClient dclient.Interface, cmResolver engineapi.ConfigmapResolver) gin.HandlerFunc {
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

		l, err := loader.New(params.Kubernetes.Version)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to initialize loader")
			return
		}

		resources, err := loader.LoadResources(l, []byte(request.Resources))
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "failed parse resources")
			return
		}

		policies, err := utils.LoadPolicies(l, []byte(request.Policies))
		if err != nil {
			c.String(http.StatusInternalServerError, "failed parse policies")
			return
		}

		config, err := loader.Load[corev1.ConfigMap](l, []byte(request.Config))
		if err != nil {
			c.String(http.StatusInternalServerError, "failed parse kyverno configmap")
			return
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
	}
}
