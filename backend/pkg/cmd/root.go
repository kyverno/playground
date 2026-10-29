package cmd

import (
	"github.com/spf13/cobra"
)

// NewCLI creates a new instance of the root CLI
func New() *cobra.Command {
	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(NewRunCommand())
	rootCmd.AddCommand(NewMCPCommand())

	return rootCmd
}
