package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigResponse struct {
	Cluster bool   `json:"cluster"`
	Sponsor string `json:"sponsor"`
}

func AddRoutes(group *gin.RouterGroup, cluster bool, sponsor string) {
	group.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, ConfigResponse{
			Cluster: cluster,
			Sponsor: sponsor,
		})
	})
}
