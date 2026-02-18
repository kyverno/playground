package authz

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/kyverno-authz/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/sdk/core"
	"github.com/kyverno/sdk/core/dispatchers"
	"github.com/kyverno/sdk/core/handlers"
	"github.com/kyverno/sdk/core/resulters"
	"github.com/kyverno/sdk/extensions/policy"
	"k8s.io/client-go/dynamic"
)

type APIConfiguration struct {
	BuiltInCrds []string
	LocalCrds   []string
}

func AddRoutes(group *gin.RouterGroup, cluster cluster.Cluster, config APIConfiguration) error {
	// handler, err := newEngineHandler(cluster, config)
	// if err != nil {
	// 	return err
	// }
	group.POST("/authz", func(ctx *gin.Context) {
		var source engine.EnvoySource
		_ = core.NewEngine(
			source,
			handlers.Handler(
				dispatchers.Sequential(
					policy.EvaluatorFactory[engine.EnvoyPolicy](),
					func(ctx context.Context, fc core.FactoryContext[engine.EnvoyPolicy, dynamic.Interface, *authv3.CheckRequest]) core.Breaker[engine.EnvoyPolicy, *authv3.CheckRequest, policy.Evaluation[*authv3.CheckResponse]] {
						return core.MakeBreakerFunc(func(_ context.Context, _ engine.EnvoyPolicy, _ *authv3.CheckRequest, out policy.Evaluation[*authv3.CheckResponse]) bool {
							return out.Result != nil
						})
					},
				),
				func(ctx context.Context, fc core.FactoryContext[engine.EnvoyPolicy, dynamic.Interface, *authv3.CheckRequest]) core.Resulter[engine.EnvoyPolicy, *authv3.CheckRequest, policy.Evaluation[*authv3.CheckResponse], policy.Evaluation[*authv3.CheckResponse]] {
					return resulters.NewFirst[engine.EnvoyPolicy, *authv3.CheckRequest](func(out policy.Evaluation[*authv3.CheckResponse]) bool {
						return out.Result != nil || out.Error != nil
					})
				},
			),
		)
	})
	return nil
}
