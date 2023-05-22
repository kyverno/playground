package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kyverno/playground/backend/pkg/engine"
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

		loader, err := NewLoader(params.Kubernetes.Version)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to initialize loader")
			return
		}

		resources, err := loader.Resources(request.Resources)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusInternalServerError, "failed parse resources")
			return
		}

		policies, err := loader.Policies(request.Policies)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed parse policies")
			return
		}

		config, err := loader.ConfigMap(request.Config)
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

func NewNamespaceHandler(client kubernetes.Interface) gin.HandlerFunc {
	if client == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, make([]string, 0))
		}
	}

	nsClient := client.CoreV1().Namespaces()

	return func(c *gin.Context) {
		list, err := nsClient.List(c, v1.ListOptions{})
		if err != nil {
			c.String(http.StatusInternalServerError, "failed  to fetch namespaces")
			return
		}

		namespaces := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			namespaces = append(namespaces, item.GetName())
		}

		c.JSON(http.StatusOK, namespaces)
	}
}

func NewResourceListHandler(dClient dclient.Interface) gin.HandlerFunc {
	if dClient == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, make([]string, 0))
		}
	}

	return func(c *gin.Context) {
		var request ListResourcesRequest
		err := c.ShouldBind(&request)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		var selector *v1.LabelSelector
		if request.Selector != nil {
			selector = &v1.LabelSelector{MatchLabels: request.Selector}
		}

		list, err := dClient.ListResource(c, request.APIVersion, request.Kind, request.Namespace, selector)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed  to fetch namespaces")
			return
		}

		resources := make([]ListResourcesResponse, 0, len(list.Items))
		for _, item := range list.Items {
			resources = append(resources, ListResourcesResponse{
				Namespace: item.GetNamespace(),
				Name:      item.GetName(),
			})
		}

		c.JSON(http.StatusOK, resources)
	}
}

func NewResourceHandler(dClient dclient.Interface) gin.HandlerFunc {
	if dClient == nil {
		return func(c *gin.Context) {
			c.JSON(http.StatusOK, make([]string, 0))
		}
	}

	return func(c *gin.Context) {
		var request GetResourceRequest
		err := c.ShouldBind(&request)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		resource, err := dClient.GetResource(c, request.APIVersion, request.Kind, request.Namespace, request.Name)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to fetch resource")
			return
		}

		// cleanup metadata
		if meta, ok := resource.Object["metadata"]; ok {
			delete(meta.(map[string]any), "managedFields")
		}

		// cleanup status
		delete(resource.Object, "status")

		// TODO: check
		// spec.strategy.rollingUpdate: Invalid value: value provided for unknown field
		if meta, ok := resource.Object["spec"]; ok {
			delete(meta.(map[string]any), "strategy")
		}

		c.JSON(http.StatusOK, resource.Object)
	}
}
