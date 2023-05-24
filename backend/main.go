package main

import (
	"github.com/kyverno/playground/backend/pkg/cmd"
)

func main() {
	// set controller-runtime logger
	// log.SetLogger(logr.Discard())
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
