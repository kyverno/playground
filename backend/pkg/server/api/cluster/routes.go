package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
)

type SearchRequest struct {
	APIVersion string            `query:"apiVersion"`
	Kind       string            `query:"kind"`
	Namespace  string            `query:"namespace"`
	Selector   map[string]string `query:"namespace"`
}

type SearchResponse = []SearchResult

type ResourceRequest struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}

func AddRoutes(group *gin.RouterGroup, cluster Cluster) error {
	group.GET("/namespaces", func(c *gin.Context) {
		namespaces, err := cluster.Namespaces(c)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to fetch namespaces")
			return
		}
		c.JSON(http.StatusOK, namespaces)
	})
	group.GET("/search", tonic.Handler(func(c *gin.Context, in *SearchRequest) (SearchResponse, error) {
		return cluster.Search(c, in.APIVersion, in.Kind, in.Namespace, in.Selector)
	}, http.StatusOK))
	group.POST("/resource", func(c *gin.Context) {
		var request ResourceRequest
		if err := c.ShouldBind(&request); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		resource, err := cluster.Get(c, request.APIVersion, request.Kind, request.Namespace, request.Name)
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
	})
	return nil
}
