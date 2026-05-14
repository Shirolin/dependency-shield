package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shiro/dependencyshield/internal/audit"
	"github.com/shiro/dependencyshield/internal/model"
	"github.com/shiro/dependencyshield/internal/prober"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit package manager configurations for security policies",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(outWriter, "🛡️ DependencyShield Audit Report")
		fmt.Fprintln(outWriter, "---------------------------------")

		tools := []struct {
			name       string
			globalPath string
			localPath  string
			auditFunc  func(string) model.AuditResult
		}{
			{"npm", prober.GetNpmrcPath(), prober.GetLocalNpmrcPath(), audit.AuditNpm},
			{"pnpm", prober.GetPnpmrcPath(), prober.GetLocalPnpmrcPath(), audit.AuditPnpm},
			{"uv", prober.GetUvConfigPath(), prober.GetLocalUvConfigPath(), audit.AuditUv},
			{"bun", prober.GetBunfigPath(), prober.GetLocalBunfigPath(), audit.AuditBun},
		}

		for _, t := range tools {
			// Global
			globalRes := t.auditFunc(t.globalPath)
			printResult(globalRes, t.name, "Global")

			// Local
			localRes := t.auditFunc(t.localPath)
			printResult(localRes, t.name, "Local")
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}

func printResult(res model.AuditResult, toolName, scope string) {
	var icon string
	var c *color.Color

	switch res.Status {
	case model.StatusPassed:
		icon = "✅"
		c = color.New(color.FgGreen)
	case model.StatusFailed:
		icon = "❌"
		c = color.New(color.FgRed)
	case model.StatusSkip:
		icon = "⚠️"
		c = color.New(color.FgYellow)
	default:
		icon = "❓"
		c = color.New(color.FgWhite)
	}

	c.Fprintf(outWriter, "[%s] %s (%s): %s", icon, toolName, scope, res.Status)
	if res.Status == model.StatusPassed {
		fmt.Fprintf(outWriter, " (Path: %s)\n", res.ConfigPath)
	} else {
		fmt.Fprintf(outWriter, " (%s)\n", res.Message)
	}
}
