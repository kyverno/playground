package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyverno/playground/backend/pkg/api"
	"github.com/kyverno/playground/backend/pkg/config"
	apicluster "github.com/kyverno/playground/backend/pkg/server/api/cluster"
	apiconfig "github.com/kyverno/playground/backend/pkg/server/api/config"
)

func AddRoutes(group *gin.RouterGroup, config config.Config, cluster apicluster.Cluster, sponsor string) error {
	if err := apiconfig.AddRoutes(group, cluster != nil, sponsor); err != nil {
		return err
	}
	if cluster != nil {
		if err := apicluster.AddRoutes(group, cluster); err != nil {
			return err
		}
	}
	if config != nil {
		if dClient, err := config.DClient(); err != nil {
			return err
		} else if cmResolver, err := config.CMResolver(); err != nil {
			return err
		} else {
			group.POST("/engine", api.NewEngineHandler(dClient, cmResolver))
			return nil
		}
	} else {
		group.POST("/engine", api.NewEngineHandler(nil, nil))
		return nil
	}
}
