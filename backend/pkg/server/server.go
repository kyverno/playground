package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kyverno/playground/backend/pkg/config"
	"github.com/kyverno/playground/backend/pkg/server/api"
	"github.com/kyverno/playground/backend/pkg/server/api/cluster"
)

//go:embed dist
var staticFiles embed.FS

const apiPrefix = "/api"

type Shutdown = func(context.Context) error

type Server interface {
	AddUIRoutes() error
	AddAPIRoutes(config.Config, cluster.Cluster, string) error
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

func (s server) AddAPIRoutes(config config.Config, cluster cluster.Cluster, sponsor string) error {
	return api.AddRoutes(s.Group(apiPrefix), config, cluster, sponsor)
}

func (s server) AddUIRoutes() error {
	fs, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		return err
	}

	s.Use(static.Serve("/", static.ServeFileSystem(&fileSystem{http.FS(fs)})))

	fileServer := http.FileServer(http.FS(fs))
	s.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, apiPrefix) {
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	return nil
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

type fileSystem struct {
	fs http.FileSystem
}

func (b *fileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *fileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}
