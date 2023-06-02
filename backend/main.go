package main

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kyverno/playground/backend/pkg/cmd"
)

func main() {
	// set controller-runtime logger
	log.SetLogger(logr.Discard())
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
