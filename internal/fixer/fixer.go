package fixer

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/shiro/dependencyshield/internal/config"
)

// FixNpmrc sets min-release-age in .npmrc
func FixNpmrc(path string) error {
	target := fmt.Sprintf("min-release-age=%s", config.NpmMinAge)
	re := regexp.MustCompile(`(?m)^min-release-age=.*$`)
	return fixFile(path, re, target, target)
}

// FixPnpmrc sets minimum-release-age in .pnpmrc
func FixPnpmrc(path string) error {
	target := fmt.Sprintf("minimum-release-age=%s", config.PnpmMinAgeMins)
	re := regexp.MustCompile(`(?m)^minimum-release-age=.*$`)
	return fixFile(path, re, target, target)
}

// FixUvConfig sets exclude-newer in uv.toml or pyproject.toml
func FixUvConfig(path string) error {
	target := fmt.Sprintf("exclude-newer = \"%s\"", config.UvExcludeNewer)
	re := regexp.MustCompile(`(?m)^exclude-newer\s*=.*$`)
	return fixFile(path, re, target, target)
}

// FixBunfig sets minimumReleaseAge in bunfig.toml
func FixBunfig(path string) error {
	targetValue := fmt.Sprintf("minimumReleaseAge = %s", config.BunMinAgeSecs)
	re := regexp.MustCompile(`(?m)^minimumReleaseAge\s*=.*$`)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(path, []byte("[install]\n"+targetValue+"\n"), 0644)
		}
		return err
	}

	sContent := string(content)
	if re.MatchString(sContent) {
		newContent := re.ReplaceAllString(sContent, targetValue)
		return os.WriteFile(path, []byte(newContent), 0644)
	}

	// Not found, try to find [install] section
	installRe := regexp.MustCompile(`(?m)^\[install\]\s*$`)
	if installRe.MatchString(sContent) {
		// Use ReplaceAllStringFunc to ensure we only replace the [install] section header once and append after it
		// Or simpler, just replace [install] with [install]\n...
		newContent := installRe.ReplaceAllString(sContent, "[install]\n"+targetValue)
		return os.WriteFile(path, []byte(newContent), 0644)
	}

	// Append at end
	if len(sContent) > 0 && !strings.HasSuffix(sContent, "\n") {
		sContent += "\n"
	}
	sContent += "[install]\n" + targetValue + "\n"
	return os.WriteFile(path, []byte(sContent), 0644)
}

func fixFile(path string, re *regexp.Regexp, target, appendStr string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(path, []byte(appendStr+"\n"), 0644)
		}
		return err
	}

	sContent := string(content)
	if re.MatchString(sContent) {
		newContent := re.ReplaceAllString(sContent, target)
		return os.WriteFile(path, []byte(newContent), 0644)
	}

	if len(sContent) > 0 && !strings.HasSuffix(sContent, "\n") {
		sContent += "\n"
	}
	sContent += appendStr + "\n"
	return os.WriteFile(path, []byte(sContent), 0644)
}
