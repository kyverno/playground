package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/api"
)

func run(sponsor, host string, port int, log bool) {
	fs, err := fs.Sub(data.StaticFiles(), "dist")
	if err != nil {
		panic(err)
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

	router.POST("/engine", api.Serve)
	router.POST("/sponsor", func(c *gin.Context) {
		c.String(http.StatusOK, sponsor)
	})

	router.StaticFS("/", http.FS(fs))
	address := fmt.Sprintf("%v:%v", host, port)
	if err := router.Run(address); err != nil {
		panic(err)
	}
}

func main() {
	host := flag.String("host", "0.0.0.0", "server host")
	port := flag.Int("port", 8080, "server port")
	mode := flag.String("mode", gin.ReleaseMode, "gin run mode")
	log := flag.Bool("log", false, "enable gin logger")
	sponsor := flag.String("sponsor", "", "sponsor text")
	flag.Parse()
	gin.SetMode(*mode)
	run(*sponsor, *host, *port, *log)
}
