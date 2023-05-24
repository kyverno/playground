package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/server/api"
	"github.com/kyverno/playground/backend/pkg/server/ui"
)

const apiPrefix = "/api"

type Shutdown = func(context.Context) error

type Server interface {
	AddUIRoutes() error
	AddAPIRoutes(cluster.Cluster, string, string) error
	Run(context.Context, string, int) Shutdown
}

type server struct {
	*gin.Engine
}

func New(enableLogger bool, enableCors bool) (Server, error) {
	router := gin.New()
	if enableLogger {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	if enableCors {
		router.Use(cors.New(cors.Config{
			AllowOrigins:  []string{"*"},
			AllowMethods:  []string{"POST", "GET", "HEAD"},
			AllowHeaders:  []string{"Origin", "Content-Type"},
			ExposeHeaders: []string{"Content-Length"},
		}))
	}
	return server{router}, nil
}

func (s server) AddAPIRoutes(cluster cluster.Cluster, sponsor string, crds string) error {
	return api.AddRoutes(s.Group(apiPrefix), cluster, sponsor, crds)
}

func (s server) AddUIRoutes() error {
	return ui.AddRoutes(s.Engine)
}

func (s server) Run(_ context.Context, host string, port int) Shutdown {
	address := fmt.Sprintf("%v:%v", host, port)
	srv := &http.Server{
		Addr:    address,
		Handler: s.Engine.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return srv.Shutdown
}
