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
		fmt.Println("🛡️ DependencyShield Audit Report")
		fmt.Println("---------------------------------")

		results := []model.AuditResult{
			audit.AuditNpm(prober.GetNpmrcPath()),
			audit.AuditPnpm(prober.GetPnpmrcPath()),
			audit.AuditUv(prober.GetUvConfigPath()),
			audit.AuditBun(prober.GetBunfigPath()),
		}

		for _, res := range results {
			printResult(res)
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}

func printResult(res model.AuditResult) {
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

	c.Printf("[%s] %s: %s", icon, res.ToolName, res.Status)
	if res.Status == model.StatusPassed {
		fmt.Printf(" (Path: %s)\n", res.ConfigPath)
	} else {
		fmt.Printf(" (%s)\n", res.Message)
	}
}
