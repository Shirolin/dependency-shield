package prober

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetNpmrcPath(t *testing.T) {
	path := GetNpmrcPath()
	if path == "" {
		t.Error("Expected path to not be empty")
	}
	if !strings.HasSuffix(path, ".npmrc") {
		t.Errorf("Expected path to end with .npmrc, got %s", path)
	}
}

func TestGetPnpmrcPath(t *testing.T) {
	path := GetPnpmrcPath()
	if path == "" {
		t.Error("Expected path to not be empty")
	}
	if !strings.HasSuffix(path, ".npmrc") {
		t.Errorf("Expected path to end with .npmrc, got %s", path)
	}
}

func TestGetUvConfigPath(t *testing.T) {
	path := GetUvConfigPath()
	if path == "" {
		t.Error("Expected path to not be empty")
	}
	if !strings.HasSuffix(path, "uv.toml") {
		t.Errorf("Expected path to end with uv.toml, got %s", path)
	}
}

func TestGetBunfigPath(t *testing.T) {
	path := GetBunfigPath()
	if path == "" {
		t.Error("Expected path to not be empty")
	}
	if !strings.HasSuffix(path, ".bunfig.toml") {
		t.Errorf("Expected path to end with .bunfig.toml, got %s", path)
	}
}

func TestFindConfigUpwards(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "prober_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	parentDir := filepath.Join(tempDir, "parent")
	childDir := filepath.Join(parentDir, "child")
	grandchildDir := filepath.Join(childDir, "grandchild")

	if err := os.MkdirAll(grandchildDir, 0755); err != nil {
		t.Fatalf("Failed to create directory structure: %v", err)
	}

	configPath := filepath.Join(parentDir, ".npmrc")
	if err := os.WriteFile(configPath, []byte("registry=https://registry.npmjs.org/"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Change working directory to grandchild
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	if err := os.Chdir(grandchildDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	foundPath := FindConfigUpwards(".npmrc")
	if foundPath == "" {
		t.Error("Expected to find .npmrc, but got empty string")
	}

	absConfigPath, _ := filepath.Abs(configPath)
	if foundPath != absConfigPath {
		t.Errorf("Expected path %s, got %s", absConfigPath, foundPath)
	}
}
