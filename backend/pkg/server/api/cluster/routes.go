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
	Selector   map[string]string `query:"selector"`
}

type SearchResponse = []SearchResult

type ResourceRequest struct {
	APIVersion string `query:"apiVersion"`
	Kind       string `query:"kind"`
	Namespace  string `query:"namespace"`
	Name       string `query:"name"`
}

func AddRoutes(group *gin.RouterGroup, cluster Cluster) error {
	group.GET("/kinds", func(c *gin.Context) {
		kinds, err := cluster.Kinds(c, "kyverno.io", "wgpolicyk8s.io")
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to fetch kinds")
			return
		}
		c.JSON(http.StatusOK, kinds)
	})
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
	group.GET("/resource", tonic.Handler(func(c *gin.Context, in *ResourceRequest) (map[string]interface{}, error) {
		resource, err := cluster.Get(c, in.APIVersion, in.Kind, in.Namespace, in.Name)
		if err != nil {
			return nil, err
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
		return resource.Object, err
	}, http.StatusOK))
	return nil
}
