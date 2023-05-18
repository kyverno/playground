package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/kyverno/playground/backend/pkg/engine"
)

type Server struct {
	k8sConfig *rest.Config
}

func (s *Server) Serve(c *gin.Context) {
	var request Request
	err := c.ShouldBind(&request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	params, err := engine.ParseParameters(request.Context)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		c.String(http.StatusBadRequest, "invalid context string")
		return
	}

	loader, err := NewLoader(params.Kubernetes.Version)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to initialize loader")
		return
	}

	resources, err := loader.Resources(request.Resources)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed parse resources")
		return
	}

	policies, err := loader.Policies(request.Policies)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed parse policies")
		return
	}

	processor, err := engine.NewProcessor(params, s.k8sConfig)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to initialize processor")
		return
	}

	results, err := processor.Run(c, policies, resources)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Writer.WriteString(err.Error()) //nolint: errcheck
		return
	}

	response := &Response{
		Policies:          policies,
		Resources:         resources,
		Mutation:          results.Mutation,
		ImageVerification: results.ImageVerification,
		Validation:        results.Validation,
		Generation:        results.Generation,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func NewServer(k8sConfig *rest.Config) *Server {
	return &Server{k8sConfig: k8sConfig}
}
