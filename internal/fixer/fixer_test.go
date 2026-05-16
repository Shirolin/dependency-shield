package fixer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shiro/dependency-shield/internal/config"
)

func TestFixNpmrc(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".npmrc")
	p := config.NewDefaultPolicy()

	// 1. Non-existent file
	err := FixNpmrc(path, p)
	if err != nil {
		t.Fatalf("FixNpmrc failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if strings.TrimSpace(string(content)) != "min-release-age=30" {
		t.Errorf("Unexpected content: %s", string(content))
	}

	// 2. Existing file but no key
	os.WriteFile(path, []byte("registry=https://registry.npmjs.org\n"), 0644)
	err = FixNpmrc(path, p)
	if err != nil {
		t.Fatalf("FixNpmrc failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "min-release-age=30") {
		t.Errorf("Key not found in content: %s", string(content))
	}

	// 3. Existing file and existing key
	os.WriteFile(path, []byte("min-release-age=0\nother=val"), 0644)
	err = FixNpmrc(path, p)
	if err != nil {
		t.Fatalf("FixNpmrc failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "min-release-age=30") || strings.Contains(string(content), "min-release-age=0") {
		t.Errorf("Replacement failed: %s", string(content))
	}

	// 4. Custom Policy (7 days)
	p7 := config.Policy{MinAgeDays: 7}
	err = FixNpmrc(path, p7)
	if err != nil {
		t.Fatalf("FixNpmrc with 7 days failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "min-release-age=7") {
		t.Errorf("Custom policy replacement failed: %s", string(content))
	}
}

func TestFixBunfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bunfig.toml")
	p := config.NewDefaultPolicy()

	// 1. Non-existent file
	err := FixBunfig(path, p)
	if err != nil {
		t.Fatalf("FixBunfig failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "[install]") || !strings.Contains(string(content), "minimumReleaseAge = 2592000") {
		t.Errorf("Unexpected content: %s", string(content))
	}

	// 2. Existing file with [install] but no key
	os.WriteFile(path, []byte("[install]\nregistry = \"https://registry.npmjs.org\"\n"), 0644)
	err = FixBunfig(path, p)
	if err != nil {
		t.Fatalf("FixBunfig failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "minimumReleaseAge = 2592000") {
		t.Errorf("Key not found in content: %s", string(content))
	}

	// 3. Existing file without [install]
	os.WriteFile(path, []byte("test = true\n"), 0644)
	err = FixBunfig(path, p)
	if err != nil {
		t.Fatalf("FixBunfig failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "[install]") || !strings.Contains(string(content), "minimumReleaseAge = 2592000") {
		t.Errorf("Section or key not found: %s", string(content))
	}

	// 4. Custom Policy (7 days)
	p7 := config.Policy{MinAgeDays: 7}
	err = FixBunfig(path, p7)
	if err != nil {
		t.Fatalf("FixBunfig with 7 days failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "minimumReleaseAge = 604800") {
		t.Errorf("Custom policy replacement failed: %s", string(content))
	}
}

func TestFixPnpmrc(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".pnpmrc")
	p := config.NewDefaultPolicy()

	err := FixPnpmrc(path, p)
	if err != nil {
		t.Fatalf("FixPnpmrc failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "minimum-release-age=43200") {
		t.Errorf("Unexpected content: %s", string(content))
	}

	// Custom Policy (7 days)
	p7 := config.Policy{MinAgeDays: 7}
	err = FixPnpmrc(path, p7)
	if err != nil {
		t.Fatalf("FixPnpmrc with 7 days failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "minimum-release-age=10080") {
		t.Errorf("Custom policy replacement failed: %s", string(content))
	}
}

func TestFixUvConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "uv.toml")
	p := config.NewDefaultPolicy()

	err := FixUvConfig(path, p)
	if err != nil {
		t.Fatalf("FixUvConfig failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "exclude-newer = \"30d\"") {
		t.Errorf("Unexpected content: %s", string(content))
	}

	// Custom Policy (7 days)
	p7 := config.Policy{MinAgeDays: 7}
	err = FixUvConfig(path, p7)
	if err != nil {
		t.Fatalf("FixUvConfig with 7 days failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "exclude-newer = \"7d\"") {
		t.Errorf("Custom policy replacement failed: %s", string(content))
	}
}
