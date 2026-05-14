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
			name   string
			path   string
			audit  func(string) model.AuditResult
			fix    func(string) error
		}{
			{"npm", prober.GetNpmrcPath(), audit.AuditNpm, fixer.FixNpmrc},
			{"pnpm", prober.GetPnpmrcPath(), audit.AuditPnpm, fixer.FixPnpmrc},
			{"uv", prober.GetUvConfigPath(), audit.AuditUv, fixer.FixUvConfig},
			{"bun", prober.GetBunfigPath(), audit.AuditBun, fixer.FixBunfig},
		}

		for _, t := range tools {
			res := t.audit(t.path)
			if res.Status == model.StatusFailed {
				err := t.fix(t.path)
				if err != nil {
					color.Red("[❌] %s: FIX FAILED (%s)\n", t.name, err.Error())
				} else {
					color.Green("[✅] %s: FIXED\n", t.name)
				}
			} else if res.Status == model.StatusPassed {
				color.Cyan("[✔] %s: ALREADY PASSED\n", t.name)
			} else {
				color.Yellow("[⚠️] %s: SKIPPED (%s)\n", t.name, res.Message)
			}
		}
	},
}

func init() {
	fixCmd.Flags().BoolVarP(&force, "force", "f", false, "Force fix without confirmation")
	rootCmd.AddCommand(fixCmd)
}
