package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

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

func New(log bool, kubeConfig string, sponsor string) (Server, error) {
	fs, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		return nil, err
	}
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
	var k8sConfig *rest.Config
	if kubeConfig != "" {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
	}

	resolver, err := config.NewResolver(k8sConfig)
	if err != nil {
		return nil, err
	}

	dClient, err := resolver.DClient()
	if err != nil {
		return nil, err
	}

	cmResolver, err := resolver.CMResolver()
	if err != nil {
		return nil, err
	}

	kubeClient, err := resolver.KubeClient()
	if err != nil {
		return nil, err
	}

	router.POST("/config", api.NewConfigHandler(kubeConfig != "", sponsor))
	router.POST("/engine", api.NewEngineHandler(dClient, cmResolver))
	router.POST("/namespaces", api.NewNamespaceHandler(kubeClient))
	router.POST("/resources", api.NewResourceListHandler(dClient))
	router.POST("/resource", api.NewResourceHandler(dClient))

	router.StaticFS("/", http.FS(fs))
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
