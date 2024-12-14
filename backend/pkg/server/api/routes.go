package api

import (
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
	apicluster "github.com/kyverno/playground/backend/pkg/server/api/cluster"
	apiconfig "github.com/kyverno/playground/backend/pkg/server/api/config"
	apiengine "github.com/kyverno/playground/backend/pkg/server/api/engine"
)

const clusterPrefix = "/cluster"

type EngineConfiguration = apiengine.APIConfiguration

type APIConfiguration struct {
	EngineConfiguration
	Versions []apiconfig.Version
	Sponsor  string
}

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, config APIConfiguration) error {
	if err := apiconfig.AddRoutes(group, cluster, config.Sponsor, config.Versions); err != nil {
		return err
	}
	if err := apiengine.AddRoutes(group, cluster, config.EngineConfiguration); err != nil {
		return err
	}
	// do not register cluster routes if we don't have a cluster
	if cluster != nil {
		if err := apicluster.AddRoutes(group.Group(clusterPrefix), cluster); err != nil {
			return err
		}
	}
	return nil
}
