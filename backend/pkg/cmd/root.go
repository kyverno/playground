package cmd

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine"
	"github.com/kyverno/playground/backend/pkg/server"
	"github.com/kyverno/playground/backend/pkg/server/api"
	"github.com/kyverno/playground/backend/pkg/utils"
)

type commandFlags struct {
	serverFlags  serverFlags
	ginFlags     ginFlags
	uiFlags      uiFlags
	engineFlags  engineFlags
	clusterFlags clusterFlags
}

type serverFlags struct {
	host string
	port int
}

type ginFlags struct {
	mode        string
	log         bool
	cors        bool
	maxBodySize int
}

type uiFlags struct {
	sponsor string
}

type engineFlags struct {
	builtInCrds []string
	localCrds   []string
}

type clusterFlags struct {
	cluster             bool
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
	res.Flags().IntVar(&command.ginFlags.maxBodySize, "gin-max-body-size", 2*1024*1024, "gin max body size")
	// ui flags
	res.Flags().StringVar(&command.uiFlags.sponsor, "ui-sponsor", "", "sponsor text")
	// engine flags
	res.Flags().StringSliceVar(&command.engineFlags.builtInCrds, "engine-builtin-crds", nil, "list of enabled builtin custom resource definitions")
	res.Flags().StringSliceVar(&command.engineFlags.localCrds, "engine-local-crds", nil, "list of folders containing custom resource definitions")
	// cluster flags
	res.Flags().BoolVar(&command.clusterFlags.cluster, "cluster", false, "enable cluster connected mode")
	clientcmd.BindOverrideFlags(&command.clusterFlags.kubeConfigOverrides, res.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return res
}

func (c *commandFlags) Run(_ *cobra.Command, _ []string) error {
	// initialise gin framework
	gin.SetMode(c.ginFlags.mode)
	tonic.SetBindHook(tonic.DefaultBindingHookMaxBodyBytes(int64(c.ginFlags.maxBodySize)))
	tonic.SetErrorHook(func(c *gin.Context, err error) (int, interface{}) {
		switch e := err.(type) {
		case engine.PolicyViolationError:
			return http.StatusBadRequest, gin.H{
				"violations": e.Violations,
				"error":      e.Error(),
				"reason":     "POLICY_VALIDATION",
			}
		default:
			return http.StatusBadRequest, gin.H{
				"error":  e.Error(),
				"reason": "ERROR",
			}
		}
	})
	// create server
	server, err := server.New(c.ginFlags.log, c.ginFlags.cors)
	if err != nil {
		return err
	}
	// register UI routes
	if err := server.AddUIRoutes(); err != nil {
		return err
	}
	apiConfig := api.APIConfiguration{
		Sponsor: c.uiFlags.sponsor,
		EngineConfiguration: api.EngineConfiguration{
			BuiltInCrds: c.engineFlags.builtInCrds,
			LocalCrds:   c.engineFlags.localCrds,
		},
	}
	// register API routes (with/without cluster support)
	if c.clusterFlags.cluster {
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
		if err := server.AddAPIRoutes(cluster, apiConfig); err != nil {
			return err
		}
	} else {
		// register API routes
		if err := server.AddAPIRoutes(cluster.NewFake(), apiConfig); err != nil {
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
