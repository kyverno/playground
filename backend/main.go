package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/server"
)

type options struct {
	host       string
	port       int
	mode       string
	log        bool
	kubeConfig string
	sponsor    string
}

func getOptions() options {
	var options options
	flag.StringVar(&options.host, "host", "0.0.0.0", "server host")
	flag.IntVar(&options.port, "port", 8080, "server port")
	flag.StringVar(&options.mode, "mode", gin.ReleaseMode, "gin run mode")
	flag.BoolVar(&options.log, "log", false, "enable gin logger")
	flag.StringVar(&options.kubeConfig, "kubeconfig", "", "enable gin logger")
	flag.StringVar(&options.sponsor, "sponsor", "", "sponsor text")
	flag.Parse()
	return options
}

func main() {
	options := getOptions()
	gin.SetMode(options.mode)
	if server, err := server.New(options.log, options.kubeConfig, options.sponsor); err != nil {
		panic(err)
	} else if err := server.Run(options.host, options.port); err != nil {
		panic(err)
	}
}
