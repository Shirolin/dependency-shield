package prober

import (
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
