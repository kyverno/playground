package engine

import (
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
)

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, crds string) error {
	handler, err := newEngineHandler(cluster, crds)
	if err != nil {
		return err
	}
	group.POST("/engine", handler)
	return nil
}
