package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/playground/backend/pkg/server"
	"github.com/kyverno/playground/backend/pkg/server/api/cluster"
	"github.com/kyverno/playground/backend/pkg/utils"
)

type options struct {
	host       string
	port       int
	mode       string
	log        bool
	cors       bool
	kubeConfig string
	sponsor    string
}

func getOptions() options {
	var options options
	flag.StringVar(&options.host, "host", "0.0.0.0", "server host")
	flag.IntVar(&options.port, "port", 8080, "server port")
	flag.StringVar(&options.mode, "mode", gin.ReleaseMode, "gin run mode")
	flag.BoolVar(&options.log, "log", false, "enable gin logger")
	flag.BoolVar(&options.cors, "cors", true, "enable gin cors")
	flag.StringVar(&options.kubeConfig, "kubeconfig", "", "enable gin logger")
	flag.StringVar(&options.sponsor, "sponsor", "", "sponsor text")
	flag.Parse()
	return options
}

func main() {
	// get options from command line parameters
	options := getOptions()
	// initialise gin framework
	gin.SetMode(options.mode)
	// create server
	server, err := server.New(options.log, options.cors)
	if err != nil {
		panic(err)
	}
	// register UI routes
	if err := server.AddUIRoutes(); err != nil {
		panic(err)
	}
	// register API routes (with/without cluster support)
	if options.kubeConfig != "" {
		// create rest config
		restConfig, err := utils.RestConfig(options.kubeConfig)
		if err != nil {
			panic(err)
		}
		// create cluster
		cluster, err := cluster.New(restConfig)
		if err != nil {
			panic(err)
		}
		// register API routes
		if err := server.AddAPIRoutes(cluster, options.sponsor); err != nil {
			panic(err)
		}
	} else {
		// register API routes
		if err := server.AddAPIRoutes(nil, options.sponsor); err != nil {
			panic(err)
		}
	}
	// run server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	shutdown := server.Run(ctx, options.host, options.port)
	<-ctx.Done()
	stop()
	if shutdown != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			panic(err)
		}
	}
}
