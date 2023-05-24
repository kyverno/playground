package main

import (
	"fmt"
	"os"

	"github.com/kyverno/playground/backend/pkg/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
