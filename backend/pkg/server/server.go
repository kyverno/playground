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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	apiServer := api.NewServer(k8sConfig)
	router.POST("/engine", apiServer.Serve)
	router.POST("/sponsor", func(c *gin.Context) {
		c.String(http.StatusOK, sponsor)
	})
	ui := router.Group("/ui")
	ui.StaticFS("/", http.FS(fs))
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
