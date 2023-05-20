package server

import (
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

type Server interface {
	Run(string, int) error
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
	router.StaticFS("/", http.FS(fs))
	return server{router}, nil
}

func (s server) Run(host string, port int) error {
	address := fmt.Sprintf("%v:%v", host, port)
	return s.Engine.Run(address)
}
