package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shield",
	Short: "DependencyShield - Security audit and fixer for package managers",
	Long:  `DependencyShield (shield) is a security CLI tool to audit and fix package manager configurations to prevent dependency-based attacks.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
