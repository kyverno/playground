package engine

import (
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
)

type APIConfiguration struct {
	BuiltInCrds []string
	LocalCrds   []string
}

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, config APIConfiguration) error {
	handler, err := newEngineHandler(cluster, config)
	if err != nil {
		return err
	}
	group.POST("/engine", handler)
	return nil
}
