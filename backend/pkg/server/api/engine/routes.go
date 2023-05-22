package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	"github.com/kyverno/kyverno/pkg/engine/context/resolvers"
	"github.com/kyverno/playground/backend/pkg/api"
	"k8s.io/client-go/kubernetes"
)

func AddRoutes(group *gin.RouterGroup, kubeClient kubernetes.Interface, dClient dclient.Interface) error {
	if kubeClient != nil {
		cmResolver, err := resolvers.NewClientBasedResolver(kubeClient)
		if err != nil {
			return err
		}
		group.POST("/engine", api.NewEngineHandler(dClient, cmResolver))
	} else {
		group.POST("/engine", api.NewEngineHandler(dClient, nil))
	}
	return nil
}
