package cel

import (
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
)

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster) error {
	handler, err := newHandler(cluster)
	if err != nil {
		return err
	}
	group.POST("/cel", handler)
	return nil
}
