package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/kyverno/playground/backend/pkg/config"
	"github.com/kyverno/playground/backend/pkg/mcp"
)

type mcpFlags struct {
	serverFlags serverFlags
	configFile  string
}

func NewMCPCommand() *cobra.Command {
	command := mcpFlags{}
	res := &cobra.Command{
		RunE: command.Run,
		Use:  "mcp",
	}
	res.Flags().StringVar(&command.configFile, "config", "", "path to an optional config file")
	// server flags
	res.Flags().StringVar(&command.serverFlags.host, "server-host", "0.0.0.0", "server host")
	res.Flags().IntVar(&command.serverFlags.port, "server-port", 9090, "server port")

	return res
}

func (c *mcpFlags) Run(_ *cobra.Command, _ []string) error {
	cfg := &config.Config{}
	err := config.Load(cfg, c.configFile)
	if err != nil {
		return err
	}

	go func() {
		s := http.NewServeMux()
		s.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", c.serverFlags.host, 8085), s); err != nil {
			fmt.Println("Failed to start health server:", err)
		}
	}()

	return mcp.New().Start(fmt.Sprintf("%s:%d", c.serverFlags.host, c.serverFlags.port))
}
