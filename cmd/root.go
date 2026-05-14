package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dependencyshield",
	Short: "A security CLI tool to audit package manager configurations",
	Long:  `DependencyShield is a security CLI tool to audit package manager configurations.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
