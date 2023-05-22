package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kyverno/playground/backend/pkg/api"
	"github.com/kyverno/playground/backend/pkg/config"
	apiconfig "github.com/kyverno/playground/backend/pkg/server/api/config"
)

func AddRoutes(group *gin.RouterGroup, config config.Config, sponsor string) error {
	if dClient, err := config.DClient(); err != nil {
		return err
	} else if cmResolver, err := config.CMResolver(); err != nil {
		return err
	} else if kubeClient, err := config.KubeClient(); err != nil {
		return err
	} else {
		if kubeClient != nil {
			group.GET("/namespaces", api.NewNamespaceHandler(kubeClient))
		}
		if dClient != nil {
			group.POST("/resources", api.NewResourceListHandler(dClient))
			group.POST("/resource", api.NewResourceHandler(dClient))
		}
		apiconfig.AddRoutes(group, kubeClient != nil, sponsor)
		group.POST("/engine", api.NewEngineHandler(dClient, cmResolver))
		return nil
	}
}
