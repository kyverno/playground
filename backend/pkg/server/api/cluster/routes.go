package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Namespace  string            `json:"namespace"`
	Selector   map[string]string `json:"selector"`
}

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
	group.POST("/resources", func(c *gin.Context) {
		var request SearchRequest
		if err := c.ShouldBind(&request); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		list, err := cluster.Search(c, request.APIVersion, request.Kind, request.Namespace, request.Selector)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to fetch resources")
			return
		}
		c.JSON(http.StatusOK, list)
	})
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
