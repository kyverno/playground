package engine

import (
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
)

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, builtInCrds ...string) error {
	handler, err := newEngineHandler(cluster, builtInCrds...)
	if err != nil {
		return err
	}
	group.POST("/engine", handler)
	return nil
}
