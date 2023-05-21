package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/api"
	"github.com/kyverno/playground/backend/pkg/config"
)

//go:embed dist
var staticFiles embed.FS

type Shutdown = func(context.Context) error

type Server interface {
	Run(context.Context, string, int) Shutdown
}

type server struct {
	*gin.Engine
}

func New(config config.Config, log bool, sponsor string) (Server, error) {
	router := gin.New()
	if log {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	apiGroup := router.Group("/api")
	if err := addApiRoutes(apiGroup, config, sponsor); err != nil {
		return nil, err
	}
	uiGroup := router.Group("/")
	if err := addUiRoutes(uiGroup); err != nil {
		return nil, err
	}
	return server{router}, nil
}

func (s server) Run(ctx context.Context, host string, port int) Shutdown {
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

func addApiRoutes(group *gin.RouterGroup, config config.Config, sponsor string) error {
	if dClient, err := config.DClient(); err != nil {
		return err
	} else if cmResolver, err := config.CMResolver(); err != nil {
		return err
	} else if kubeClient, err := config.KubeClient(); err != nil {
		return err
	} else {
		if kubeClient != nil {
			group.POST("/namespaces", api.NewNamespaceHandler(kubeClient))
		}
		if dClient != nil {
			group.POST("/resources", api.NewResourceListHandler(dClient))
			group.POST("/resource", api.NewResourceHandler(dClient))
		}
		group.POST("/config", api.NewConfigHandler(kubeClient != nil, sponsor))
		group.POST("/engine", api.NewEngineHandler(dClient, cmResolver))
		return nil
	}
}

func addUiRoutes(group *gin.RouterGroup) error {
	fs, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		return err
	}
	group.StaticFS("/", http.FS(fs))
	return nil
}
