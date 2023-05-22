package api

import (
	"github.com/gin-gonic/gin"
	apicluster "github.com/kyverno/playground/backend/pkg/server/api/cluster"
	apiconfig "github.com/kyverno/playground/backend/pkg/server/api/config"
	apiengine "github.com/kyverno/playground/backend/pkg/server/api/engine"
)

const clusterPrefix = "/cluster"

func AddRoutes(group *gin.RouterGroup, cluster apicluster.Cluster, sponsor string) error {
	if err := apiconfig.AddRoutes(group, cluster != nil, sponsor); err != nil {
		return err
	}
	if cluster != nil {
		if err := apicluster.AddRoutes(group.Group(clusterPrefix), cluster); err != nil {
			return err
		}
		if err := apiengine.AddRoutes(group, cluster.KubeClient(), cluster.DClient()); err != nil {
			return err
		}
	} else {
		if err := apiengine.AddRoutes(group, nil, nil); err != nil {
			return err
		}
	}
	return nil
}
