package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shiro/dependency-shield/internal/audit"
	"github.com/shiro/dependency-shield/internal/config"
	"github.com/shiro/dependency-shield/internal/model"
	"github.com/shiro/dependency-shield/internal/prober"
	"github.com/spf13/cobra"
)

var minAgeDays int
var env string

var auditCmd = &cobra.Command{
	Use:   "audit [paths...]",
	Short: "Audit package manager configurations for security policies",
	Long:  "Audit package manager configurations for security policies. If paths are provided, only those paths (and subdirectories) are scanned. Otherwise, global and current directory configurations are scanned.",
	Run: func(cmd *cobra.Command, args []string) {
		// Override minAgeDays based on environment preset
		switch env {
		case config.EnvLocal:
			minAgeDays = config.EnvLocalDays
		case config.EnvCI:
			minAgeDays = config.EnvCIDays
		case config.EnvProd:
			minAgeDays = config.EnvProdDays
		case "":
			// Keep original (default or via --days)
		default:
			fmt.Fprintf(outWriter, "⚠️  Unknown environment '%s', falling back to %d days\n", env, minAgeDays)
		}

		fmt.Fprintln(outWriter, "🛡️ DependencyShield Audit Report")
		fmt.Fprintf(outWriter, "Policy: Minimum release age = %d days\n", minAgeDays)
		fmt.Fprintln(outWriter, "---------------------------------")

		policy := config.Policy{MinAgeDays: minAgeDays}

		tools := []struct {
			name       string
			globalPath string
			localPaths []string
			auditFunc  func(string, config.Policy) model.AuditResult
		}{
			{"npm", prober.GetNpmrcPath(), prober.GetLocalNpmrcPaths(), audit.AuditNpm},
			{"pnpm", prober.GetPnpmrcPath(), prober.GetLocalPnpmrcPaths(), audit.AuditPnpm},
			{"uv", prober.GetUvConfigPath(), prober.GetLocalUvConfigPaths(), audit.AuditUv},
			{"bun", prober.GetBunfigPath(), prober.GetLocalBunfigPaths(), audit.AuditBun},
		}

		// If arguments provided, override paths
		if len(args) > 0 {
			for i := range tools {
				tools[i].globalPath = "" // Disable default global
				tools[i].localPaths = nil // Clear default locals
				for _, arg := range args {
					foundPaths := prober.GetConfigsInPath(tools[i].name, arg)
					tools[i].localPaths = append(tools[i].localPaths, foundPaths...)
				}
			}
		}

		for _, t := range tools {
			// Only audit if the tool is installed (or if we have explicit paths)
			if !prober.IsToolInstalled(t.name) && len(t.localPaths) == 0 {
				continue
			}

			// Global
			if t.globalPath != "" {
				globalRes := t.auditFunc(t.globalPath, policy)
				printResult(globalRes, t.name, "Global")
			}

			// Local(s) / Specified
			scope := "Local"
			if len(args) > 0 {
				scope = "Target"
			}

			if len(t.localPaths) == 0 && t.globalPath != "" {
				printResult(model.AuditResult{Status: model.StatusSkip, Message: "No local config found"}, t.name, scope)
			} else {
				for _, lp := range t.localPaths {
					localRes := t.auditFunc(lp, policy)
					printResult(localRes, t.name, scope)
				}
			}
		}

		// 输出建议
		fmt.Fprintln(outWriter, "\n💡 Security Recommendations:")
		if minAgeDays < 15 {
			fmt.Fprintln(outWriter, "- [Local] Your policy is set to Local development (7-14 days). This balances speed and safety.")
		}
		if minAgeDays >= 15 && minAgeDays < 30 {
			fmt.Fprintln(outWriter, "- [CI/Test] Your policy is set to CI/Test (15 days). Recommended for shared staging environments.")
		}
		if minAgeDays >= 30 {
			fmt.Fprintln(outWriter, "- [Prod] Your policy is set to Production (30 days). Maximum protection against supply chain attacks.")
		}

		fmt.Fprintln(outWriter, "- Suggestion: Use '--env' to quickly switch presets: 'audit -e local', 'audit -e ci', or 'audit -e prod'.")
	},
}

func init() {
	auditCmd.Flags().IntVarP(&minAgeDays, "days", "d", config.DefaultMinAgeDays, "Minimum release age in days")
	auditCmd.Flags().StringVarP(&env, "env", "e", "", "Environment preset: local (7d), ci (15d), prod (30d). Overrides --days.")
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
	} else if res.Status == model.StatusFailed {
		fmt.Fprintf(outWriter, " (%s, Path: %s)\n", res.Message, res.ConfigPath)
	} else {
		fmt.Fprintf(outWriter, " (%s)\n", res.Message)
	}
}
