package fixer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixNpmrc(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".npmrc")

	// 1. Non-existent file
	err := FixNpmrc(path)
	if err != nil {
		t.Fatalf("FixNpmrc failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if strings.TrimSpace(string(content)) != "min-release-age=30" {
		t.Errorf("Unexpected content: %s", string(content))
	}

	// 2. Existing file but no key
	os.WriteFile(path, []byte("registry=https://registry.npmjs.org\n"), 0644)
	err = FixNpmrc(path)
	if err != nil {
		t.Fatalf("FixNpmrc failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "min-release-age=30") {
		t.Errorf("Key not found in content: %s", string(content))
	}

	// 3. Existing file and existing key
	os.WriteFile(path, []byte("min-release-age=0\nother=val"), 0644)
	err = FixNpmrc(path)
	if err != nil {
		t.Fatalf("FixNpmrc failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "min-release-age=30") || strings.Contains(string(content), "min-release-age=0") {
		t.Errorf("Replacement failed: %s", string(content))
	}
}

func TestFixBunfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bunfig.toml")

	// 1. Non-existent file
	err := FixBunfig(path)
	if err != nil {
		t.Fatalf("FixBunfig failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "[install]") || !strings.Contains(string(content), "minimumReleaseAge = 2592000") {
		t.Errorf("Unexpected content: %s", string(content))
	}

	// 2. Existing file with [install] but no key
	os.WriteFile(path, []byte("[install]\nregistry = \"https://registry.npmjs.org\"\n"), 0644)
	err = FixBunfig(path)
	if err != nil {
		t.Fatalf("FixBunfig failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "minimumReleaseAge = 2592000") {
		t.Errorf("Key not found in content: %s", string(content))
	}

	// 3. Existing file without [install]
	os.WriteFile(path, []byte("test = true\n"), 0644)
	err = FixBunfig(path)
	if err != nil {
		t.Fatalf("FixBunfig failed: %v", err)
	}
	content, _ = os.ReadFile(path)
	if !strings.Contains(string(content), "[install]") || !strings.Contains(string(content), "minimumReleaseAge = 2592000") {
		t.Errorf("Section or key not found: %s", string(content))
	}
}

func TestFixPnpmrc(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".pnpmrc")

	err := FixPnpmrc(path)
	if err != nil {
		t.Fatalf("FixPnpmrc failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "minimum-release-age=43200") {
		t.Errorf("Unexpected content: %s", string(content))
	}
}

func TestFixUvConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "uv.toml")

	err := FixUvConfig(path)
	if err != nil {
		t.Fatalf("FixUvConfig failed: %v", err)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), "exclude-newer = \"30d\"") {
		t.Errorf("Unexpected content: %s", string(content))
	}
}
