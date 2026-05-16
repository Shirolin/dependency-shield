package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var outWriter io.Writer = os.Stdout
var inReader io.Reader = os.Stdin

func SetOut(w io.Writer) {
	outWriter = w
}

func SetIn(r io.Reader) {
	inReader = r
}

var rootCmd = &cobra.Command{
	Use:     "shield",
	Version: "v1.0.0",
	Short:   "DependencyShield - Security audit and fixer for package managers",
	Long:    `DependencyShield (shield) is a security CLI tool to audit and fix package manager configurations to prevent dependency-based attacks.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
