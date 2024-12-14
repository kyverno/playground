package config

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
)

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ConfigResponse struct {
	Cluster  bool      `json:"cluster"`
	Sponsor  string    `json:"sponsor"`
	Versions []Version `json:"versions"`
}

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, sponsor string, versions []Version) error {
	group.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, ConfigResponse{
			Cluster:  cluster != nil && !cluster.IsFake(),
			Sponsor:  sponsor,
			Versions: versions,
		})
	})
	return nil
}
