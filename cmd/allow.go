package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/shiro/dependency-shield/internal/prober"
	"github.com/spf13/cobra"
)

var allowCmd = &cobra.Command{
	Use:   "allow [package-names...]",
	Short: "Add packages to the security policy whitelist (pnpm/uv)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(outWriter, "🛡️ DependencyShield Whitelist Manager")
		fmt.Fprintln(outWriter, "---------------------------------")

		if len(args) > 0 {
			// Direct mode
			handleDirectAllow(args)
		} else {
			// Interactive mode
			handleInteractiveAllow()
		}
	},
}

func handleDirectAllow(args []string) {
	foundAny := false

	// 1. Check pnpm (.npmrc)
	npmrcPaths := prober.GetLocalNpmrcPaths()
	if globalNpmrc := prober.GetNpmrcPath(); globalNpmrc != "" {
		npmrcPaths = append(npmrcPaths, globalNpmrc)
	}

	for _, p := range npmrcPaths {
		if _, err := os.Stat(p); err == nil {
			if err := allowPnpm(p, args); err == nil {
				color.New(color.FgGreen).Fprintf(outWriter, "[✅] pnpm: Added %v to whitelist at %s\n", args, p)
				foundAny = true
			}
		}
	}

	// 2. Check uv (uv.toml / pyproject.toml)
	uvPaths := prober.GetLocalUvConfigPaths()
	if globalUv := prober.GetUvConfigPath(); globalUv != "" {
		uvPaths = append(uvPaths, globalUv)
	}

	for _, p := range uvPaths {
		if _, err := os.Stat(p); err == nil {
			if err := allowUv(p, args); err == nil {
				color.New(color.FgGreen).Fprintf(outWriter, "[✅] uv: Added %v to whitelist at %s\n", args, p)
				foundAny = true
			}
		}
	}

	if !foundAny {
		color.New(color.FgYellow).Fprintln(outWriter, "⚠️ No supported configuration files found to update.")
	}
}

type pkgItem struct {
	name       string
	tool       string // "pnpm", "uv"
	isExempted bool
	configPath string
}

func handleInteractiveAllow() {
	allPkgs := collectAllDependencies()
	if len(allPkgs) == 0 {
		fmt.Fprintln(outWriter, "No dependencies found in current project.")
		return
	}

	// Sort by name
	sort.Slice(allPkgs, func(i, j int) bool {
		return allPkgs[i].name < allPkgs[j].name
	})

	scanner := bufio.NewScanner(inReader)
	for {
		fmt.Fprintln(outWriter, "\nCurrent Dependencies & Exemption Status:")
		fmt.Fprintln(outWriter, "-----------------------------------------")
		for i, p := range allPkgs {
			status := "[ ]"
			policy := fmt.Sprintf("%d days", minAgeDays) // Use the global minAgeDays if available
			if p.isExempted {
				status = color.GreenString("[x]")
				policy = color.CyanString("0 days (EXEMPTED)")
			}
			fmt.Fprintf(outWriter, "%2d. %s %-30s (Policy: %s) [%s]\n", i+1, status, p.name, policy, p.tool)
		}
		fmt.Fprintln(outWriter, "\nCommands: <number> to toggle, 'q' to quit, 's' to save and exit")
		fmt.Fprintf(outWriter, "Action: ")

		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if input == "q" {
			fmt.Fprintln(outWriter, "Cancelled.")
			return
		}
		if input == "s" {
			saveChanges(allPkgs)
			fmt.Fprintln(outWriter, "Changes saved.")
			return
		}

		idx, err := strconv.Atoi(input)
		if err == nil && idx > 0 && idx <= len(allPkgs) {
			allPkgs[idx-1].isExempted = !allPkgs[idx-1].isExempted
		} else {
			fmt.Fprintln(outWriter, color.RedString("Invalid input."))
		}
	}
}

func collectAllDependencies() []pkgItem {
	var items []pkgItem

	// Node.js (package.json)
	wd, _ := os.Getwd()
	pkgJsonPath := filepath.Join(wd, "package.json")
	if data, err := os.ReadFile(pkgJsonPath); err == nil {
		var pkg map[string]interface{}
		if err := json.Unmarshal(data, &pkg); err == nil {
			deps, _ := pkg["dependencies"].(map[string]interface{})
			devDeps, _ := pkg["devDependencies"].(map[string]interface{})
			
			// Load whitelist for check
			whitelist := getPnpmWhitelist()

			for name := range deps {
				items = append(items, pkgItem{name: name, tool: "pnpm", isExempted: whitelist[name]})
			}
			for name := range devDeps {
				items = append(items, pkgItem{name: name, tool: "pnpm", isExempted: whitelist[name]})
			}
		}
	}

	// Python (uv.toml / pyproject.toml is skipped for brevity but can be added similarly)
	
	return items
}

func getPnpmWhitelist() map[string]bool {
	whitelist := make(map[string]bool)
	paths := prober.GetLocalNpmrcPaths()
	if g := prober.GetNpmrcPath(); g != "" {
		paths = append(paths, g)
	}

	for _, p := range paths {
		if content, err := os.ReadFile(p); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "minimum-release-age-exclude[]=") {
					parts := strings.SplitN(trimmed, "=", 2)
					if len(parts) == 2 {
						whitelist[strings.TrimSpace(parts[1])] = true
					}
				}
			}
		}
	}
	return whitelist
}

func saveChanges(items []pkgItem) {
	// Group by tool and identify what needs to be added/removed
	pnpmToExempt := []string{}
	pnpmToClear := []string{}

	for _, item := range items {
		if item.tool == "pnpm" {
			if item.isExempted {
				pnpmToExempt = append(pnpmToExempt, item.name)
			} else {
				pnpmToClear = append(pnpmToClear, item.name)
			}
		}
	}

	// Update .npmrc
	paths := prober.GetLocalNpmrcPaths()
	if g := prober.GetNpmrcPath(); g != "" {
		paths = append(paths, g)
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			updatePnpmExemptions(p, pnpmToExempt, pnpmToClear)
		}
	}
}

func updatePnpmExemptions(path string, toExempt, toClear []string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	newLines := []string{}
	
	clearMap := make(map[string]bool)
	for _, c := range toClear {
		clearMap[c] = true
	}
	
	existing := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "minimum-release-age-exclude[]=") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				pkg := strings.TrimSpace(parts[1])
				if clearMap[pkg] {
					continue // Remove it
				}
				existing[pkg] = true
			}
		}
		newLines = append(newLines, line)
	}

	for _, pkg := range toExempt {
		if !existing[pkg] {
			newLines = append(newLines, fmt.Sprintf("minimum-release-age-exclude[]=%s", pkg))
		}
	}

	os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
}

func allowPnpm(path string, pkgs []string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0, len(lines)+len(pkgs))
	
	existing := make(map[string]bool)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "minimum-release-age-exclude[]=") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				existing[strings.TrimSpace(parts[1])] = true
			}
		}
		newLines = append(newLines, line)
	}

	added := false
	for _, pkg := range pkgs {
		if !existing[pkg] {
			newLines = append(newLines, fmt.Sprintf("minimum-release-age-exclude[]=%s", pkg))
			added = true
		}
	}

	if !added {
		return nil // Nothing to add
	}

	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
}

func allowUv(path string, pkgs []string) error {
	// Simple implementation: just append to the file if it's uv.toml
	// For pyproject.toml it's more complex, but we'll try a simple approach first
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	sContent := string(content)
	added := false
	
	for _, pkg := range pkgs {
		if !strings.Contains(sContent, fmt.Sprintf("%s = false", pkg)) {
			if !strings.Contains(sContent, "[tool.uv.exclude-newer-package]") && !strings.Contains(sContent, "exclude-newer-package = {") {
				if strings.HasSuffix(sContent, "\n") {
					sContent += "\n[tool.uv.exclude-newer-package]\n"
				} else {
					sContent += "\n\n[tool.uv.exclude-newer-package]\n"
				}
			}
			sContent += fmt.Sprintf("%s = false\n", pkg)
			added = true
		}
	}

	if !added {
		return nil
	}

	return os.WriteFile(path, []byte(sContent), 0644)
}

func init() {
	rootCmd.AddCommand(allowCmd)
}
