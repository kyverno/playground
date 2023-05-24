package cmd

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kyverno/playground/backend/pkg/server"
	"github.com/kyverno/playground/backend/pkg/server/api/cluster"
	"github.com/kyverno/playground/backend/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

type commandFlags struct {
	serverFlags  serverFlags
	ginFlags     ginFlags
	configFlags  configFlags
	clusterFlags clusterFlags
}

type serverFlags struct {
	host string
	port int
}

type ginFlags struct {
	mode string
	log  bool
	cors bool
}

type configFlags struct {
	cluster bool
	sponsor string
}

type clusterFlags struct {
	kubeConfigOverrides clientcmd.ConfigOverrides
}

func NewRootCommand() *cobra.Command {
	command := commandFlags{}
	res := &cobra.Command{
		RunE: command.Run,
	}
	// server flags
	res.Flags().StringVar(&command.serverFlags.host, "server-host", "0.0.0.0", "server host")
	res.Flags().IntVar(&command.serverFlags.port, "server-port", 8080, "server port")
	// gin flags
	res.Flags().StringVar(&command.ginFlags.mode, "gin-mode", gin.ReleaseMode, "gin run mode")
	res.Flags().BoolVar(&command.ginFlags.log, "gin-log", false, "enable gin logger")
	res.Flags().BoolVar(&command.ginFlags.cors, "gin-cors", true, "enable gin cors")
	// config flags
	res.Flags().StringVar(&command.configFlags.sponsor, "sponsor", "", "sponsor text")
	res.Flags().BoolVar(&command.configFlags.cluster, "cluster", false, "enable cluster connected mode")
	// cluster flags
	clientcmd.BindOverrideFlags(&command.clusterFlags.kubeConfigOverrides, res.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return res
}

func (c *commandFlags) Run(_ *cobra.Command, args []string) error {
	// initialise gin framework
	gin.SetMode(c.ginFlags.mode)
	// create server
	server, err := server.New(c.ginFlags.log, c.ginFlags.cors)
	if err != nil {
		return err
	}
	// register UI routes
	if err := server.AddUIRoutes(); err != nil {
		return err
	}
	// register API routes (with/without cluster support)
	if c.configFlags.cluster {
		// create rest config
		restConfig, err := utils.RestConfig(c.clusterFlags.kubeConfigOverrides)
		if err != nil {
			return err
		}
		// create cluster
		cluster, err := cluster.New(restConfig)
		if err != nil {
			return err
		}
		// register API routes
		if err := server.AddAPIRoutes(cluster, c.configFlags.sponsor); err != nil {
			return err
		}
	} else {
		// register API routes
		if err := server.AddAPIRoutes(nil, c.configFlags.sponsor); err != nil {
			return err
		}
	}
	// run server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	shutdown := server.Run(ctx, c.serverFlags.host, c.serverFlags.port)
	<-ctx.Done()
	stop()
	if shutdown != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
