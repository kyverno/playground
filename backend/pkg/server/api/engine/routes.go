package engine

import (
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
)

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, config crd.APIConfiguration) error {
	handler, err := newEngineHandler(cluster, config)
	if err != nil {
		return err
	}
	group.POST("/engine", handler)
	return nil
}
