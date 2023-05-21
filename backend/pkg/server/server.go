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

const apiPrefix = "/api"
const uiPrefix = "/ui"

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
	apiGroup := router.Group(apiPrefix)
	if err := addAPIRoutes(apiGroup, config, sponsor); err != nil {
		return nil, err
	}
	uiGroup := router.Group(uiPrefix)
	if err := addUIRoutes(uiGroup); err != nil {
		return nil, err
	}
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, uiPrefix)
	})
	return server{router}, nil
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

func addAPIRoutes(group *gin.RouterGroup, config config.Config, sponsor string) error {
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
		group.GET("/config", api.NewConfigHandler(kubeClient != nil, sponsor))
		group.POST("/engine", api.NewEngineHandler(dClient, cmResolver))
		return nil
	}
}

func addUIRoutes(group *gin.RouterGroup) error {
	fs, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		return err
	}
	group.StaticFS("/", http.FS(fs))
	return nil
}
