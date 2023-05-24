package engine

import (
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/api"
	"github.com/kyverno/playground/backend/pkg/cluster"
)

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster) error {
	handler, err := api.NewEngineHandler(cluster)
	if err != nil {
		return err
	}
	group.POST("/engine", handler)
	return nil
}
