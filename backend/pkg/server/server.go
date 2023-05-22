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
)

//go:embed dist
var staticFiles embed.FS

const apiPrefix = "/api"

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
	if err := api.AddRoutes(router.Group(apiPrefix), config, sponsor); err != nil {
		return nil, err
	}
	if err := addUIRoutes(router); err != nil {
		return nil, err
	}
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

func addUIRoutes(router *gin.Engine) error {
	fs, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		return err
	}

	router.Use(static.Serve("/", static.ServeFileSystem(&fileSystem{http.FS(fs)})))

	fileServer := http.FileServer(http.FS(fs))
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, apiPrefix) {
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	return nil
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
