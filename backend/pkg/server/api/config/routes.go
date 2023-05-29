package config

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
)

type ConfigResponse struct {
	Cluster bool   `json:"cluster"`
	Sponsor string `json:"sponsor"`
}

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, sponsor string) error {
	group.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, ConfigResponse{
			Cluster: cluster != nil && !cluster.IsFake(),
			Sponsor: sponsor,
		})
	})
	return nil
}
