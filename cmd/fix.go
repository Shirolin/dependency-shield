package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shiro/dependency-shield/internal/audit"
	"github.com/shiro/dependency-shield/internal/config"
	"github.com/shiro/dependency-shield/internal/fixer"
	"github.com/shiro/dependency-shield/internal/model"
	"github.com/shiro/dependency-shield/internal/prober"
	"github.com/spf13/cobra"
)

var force bool

var fixMinAgeDays int

var fixCmd = &cobra.Command{
	Use:   "fix [paths...]",
	Short: "Fix security policy violations in package manager configurations",
	Long:  "Fix security policy violations in package manager configurations. If paths are provided, only those paths (and subdirectories) are scanned. Otherwise, global and current directory configurations are scanned.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(outWriter, "🛠️ DependencyShield Fix Report")
		fmt.Fprintf(outWriter, "Policy: Minimum release age = %d days\n", fixMinAgeDays)
		fmt.Fprintln(outWriter, "---------------------------------")

		policy := config.Policy{MinAgeDays: fixMinAgeDays}

		tools := []struct {
			name       string
			globalPath string
			localPaths []string
			auditFunc  func(string, config.Policy) model.AuditResult
			fixFunc    func(string, config.Policy) error
		}{
			{"npm", prober.GetNpmrcPath(), prober.GetLocalNpmrcPaths(), audit.AuditNpm, fixer.FixNpmrc},
			{"pnpm", prober.GetPnpmrcPath(), prober.GetLocalPnpmrcPaths(), audit.AuditPnpm, fixer.FixPnpmrc},
			{"uv", prober.GetUvConfigPath(), prober.GetLocalUvConfigPaths(), audit.AuditUv, fixer.FixUvConfig},
			{"bun", prober.GetBunfigPath(), prober.GetLocalBunfigPaths(), audit.AuditBun, fixer.FixBunfig},
		}

		// If arguments provided, override paths
		if len(args) > 0 {
			for i := range tools {
				tools[i].globalPath = ""  // Disable default global
				tools[i].localPaths = nil // Clear default locals
				for _, arg := range args {
					foundPaths := prober.GetConfigsInPath(tools[i].name, arg)
					tools[i].localPaths = append(tools[i].localPaths, foundPaths...)
				}
			}
		}

		for _, t := range tools {
			// Only fix if the tool is installed (or if we have explicit paths)
			if !prober.IsToolInstalled(t.name) && len(t.localPaths) == 0 {
				continue
			}

			// Fix Global
			if t.globalPath != "" {
				fixTool(t.name, "Global", t.globalPath, policy, t.auditFunc, t.fixFunc)
			}

			// Fix Local(s) / Specified
			scope := "Local"
			if len(args) > 0 {
				scope = "Target"
			}

			if len(t.localPaths) == 0 && t.globalPath != "" {
				color.New(color.FgYellow).Fprintf(outWriter, "[⚠️] %s (%s): SKIPPED (No local config found)\n", t.name, scope)
			} else {
				for _, lp := range t.localPaths {
					fixTool(t.name, scope, lp, policy, t.auditFunc, t.fixFunc)
				}
			}
		}
	},
}

func askConfirmation(message string) bool {
	if force {
		return true
	}
	color.New(color.FgYellow).Fprintf(outWriter, "%s [y/N]: ", message)
	var response string
	_, err := fmt.Fscanln(inReader, &response)
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func fixTool(name, scope, path string, policy config.Policy, auditFunc func(string, config.Policy) model.AuditResult, fixFunc func(string, config.Policy) error) {
	if path == "" && scope == "Local" {
		color.New(color.FgYellow).Fprintf(outWriter, "[⚠️] %s (%s): SKIPPED (No local config found)\n", name, scope)
		return
	}

	res := auditFunc(path, policy)
	if res.Status == model.StatusFailed {
		msg := fmt.Sprintf("Fix %s (%s) configuration at %s?", name, scope, path)
		if askConfirmation(msg) {
			err := fixFunc(path, policy)
			if err != nil {
				color.New(color.FgRed).Fprintf(outWriter, "[❌] %s (%s): FIX FAILED (%s, Path: %s)\n", name, scope, err.Error(), path)
			} else {
				color.New(color.FgGreen).Fprintf(outWriter, "[✅] %s (%s): FIXED (Path: %s)\n", name, scope, path)
			}
		} else {
			color.New(color.FgYellow).Fprintf(outWriter, "[⚠️] %s (%s): FIX CANCELLED BY USER (Path: %s)\n", name, scope, path)
		}
	} else if res.Status == model.StatusPassed {
		color.New(color.FgCyan).Fprintf(outWriter, "[✔] %s (%s): ALREADY PASSED (Path: %s)\n", name, scope, path)
	} else {
		color.New(color.FgYellow).Fprintf(outWriter, "[⚠️] %s (%s): SKIPPED (%s, Path: %s)\n", name, scope, res.Message, path)
	}
}

func init() {
	fixCmd.Flags().BoolVarP(&force, "force", "f", false, "Force fix without confirmation")
	fixCmd.Flags().IntVarP(&fixMinAgeDays, "days", "d", config.DefaultMinAgeDays, "Minimum release age in days (Recommended: 7-14 for dev, 30 for prod)")
	rootCmd.AddCommand(fixCmd)
}
