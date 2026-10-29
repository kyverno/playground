package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/playground"
)

func newEngineHandler(cl cluster.Cluster, config crd.APIConfiguration) (gin.HandlerFunc, error) {
	return tonic.Handler(func(ctx *gin.Context, in *playground.EngineRequest) (*playground.EngineResponse, error) {
		return playground.Run(ctx, cl, in, config)
	}, http.StatusOK), nil
}
