package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shiro/dependencyshield/internal/audit"
	"github.com/shiro/dependencyshield/internal/fixer"
	"github.com/shiro/dependencyshield/internal/model"
	"github.com/shiro/dependencyshield/internal/prober"
	"github.com/spf13/cobra"
)

var force bool

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix security policy violations in package manager configurations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🛠️ DependencyShield Fix Report")
		fmt.Println("---------------------------------")

		tools := []struct {
			name       string
			globalPath string
			localPath  string
			auditFunc  func(string) model.AuditResult
			fixFunc    func(string) error
		}{
			{"npm", prober.GetNpmrcPath(), prober.GetLocalNpmrcPath(), audit.AuditNpm, fixer.FixNpmrc},
			{"pnpm", prober.GetPnpmrcPath(), prober.GetLocalPnpmrcPath(), audit.AuditPnpm, fixer.FixPnpmrc},
			{"uv", prober.GetUvConfigPath(), prober.GetLocalUvConfigPath(), audit.AuditUv, fixer.FixUvConfig},
			{"bun", prober.GetBunfigPath(), prober.GetLocalBunfigPath(), audit.AuditBun, fixer.FixBunfig},
		}

		for _, t := range tools {
			// Fix Global
			fixTool(t.name, "Global", t.globalPath, t.auditFunc, t.fixFunc)

			// Fix Local
			fixTool(t.name, "Local", t.localPath, t.auditFunc, t.fixFunc)
		}
	},
}

func fixTool(name, scope, path string, auditFunc func(string) model.AuditResult, fixFunc func(string) error) {
	if path == "" && scope == "Local" {
		color.Yellow("[⚠️] %s (%s): SKIPPED (No local config found)\n", name, scope)
		return
	}

	res := auditFunc(path)
	if res.Status == model.StatusFailed {
		err := fixFunc(path)
		if err != nil {
			color.Red("[❌] %s (%s): FIX FAILED (%s)\n", name, scope, err.Error())
		} else {
			color.Green("[✅] %s (%s): FIXED\n", name, scope)
		}
	} else if res.Status == model.StatusPassed {
		color.Cyan("[✔] %s (%s): ALREADY PASSED\n", name, scope)
	} else {
		color.Yellow("[⚠️] %s (%s): SKIPPED (%s)\n", name, scope, res.Message)
	}
}

func init() {
	fixCmd.Flags().BoolVarP(&force, "force", "f", false, "Force fix without confirmation")
	rootCmd.AddCommand(fixCmd)
}
