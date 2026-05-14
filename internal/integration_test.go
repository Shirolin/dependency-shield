package internal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shiro/dependency-shield/internal/audit"
	"github.com/shiro/dependency-shield/internal/model"
	"github.com/shiro/dependency-shield/internal/prober"
)

func TestNpmConfigurationHierarchy(t *testing.T) {
	// 1. Use a temporary directory to simulate a user's home and a project directory.
	tmpDir := t.TempDir()
	mockHome := filepath.Join(tmpDir, "home")
	mockProject := filepath.Join(tmpDir, "project")

	if err := os.MkdirAll(mockHome, 0755); err != nil {
		t.Fatalf("failed to create mock home: %v", err)
	}
	if err := os.MkdirAll(mockProject, 0755); err != nil {
		t.Fatalf("failed to create mock project: %v", err)
	}

	// 2. Create a global '.npmrc' in the mock home with 'min-release-age=30' (PASSED).
	globalPath := filepath.Join(mockHome, ".npmrc")
	if err := os.WriteFile(globalPath, []byte("min-release-age=30\n"), 0644); err != nil {
		t.Fatalf("failed to write global .npmrc: %v", err)
	}

	// 3. Create a local '.npmrc' in the mock project directory with 'min-release-age=10' (FAILED).
	localPath := filepath.Join(mockProject, ".npmrc")
	if err := os.WriteFile(localPath, []byte("min-release-age=10\n"), 0644); err != nil {
		t.Fatalf("failed to write local .npmrc: %v", err)
	}

	// Mock HOME for prober to find the global one
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", mockHome)
	defer os.Setenv("HOME", originalHome)
	
	// For Windows support in prober
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", mockHome)
	defer os.Setenv("USERPROFILE", originalUserProfile)

	originalAppData := os.Getenv("APPDATA")
	os.Setenv("APPDATA", mockHome)
	defer os.Setenv("APPDATA", originalAppData)

	// For Linux/macOS XDG support
	originalXdg := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", mockHome)
	defer os.Setenv("XDG_CONFIG_HOME", originalXdg)

	// Implementation Hint: Calls 'audit.AuditNpm' on two different paths and asserts the expected statuses.
	t.Run("Direct Audit", func(t *testing.T) {
		resGlobal := audit.AuditNpm(globalPath)
		if resGlobal.Status != model.StatusPassed {
			t.Errorf("Global .npmrc should pass, got %s", resGlobal.Status)
		}

		resLocal := audit.AuditNpm(localPath)
		if resLocal.Status != model.StatusFailed {
			t.Errorf("Local .npmrc should fail, got %s", resLocal.Status)
		}
	})

	// 4. The test should attempt to audit both and expect a result that includes both the global and local findings.
	t.Run("Automated Discovery", func(t *testing.T) {
		// Change working directory to mock project so prober can find the local .npmrc
		oldWd, _ := os.Getwd()
		if err := os.Chdir(mockProject); err != nil {
			t.Fatalf("failed to change wd: %v", err)
		}
		defer os.Chdir(oldWd)

		results := []model.AuditResult{
			audit.AuditNpm(prober.GetNpmrcPath()),
			audit.AuditNpm(prober.GetLocalNpmrcPath()),
		}

		// We expect 2 results if hierarchy is supported
		if len(results) < 2 {
			t.Errorf("Expected at least 2 audit results for npm (global and local), but found %d.", len(results))
		}

		foundGlobal := false
		foundLocal := false
		for _, res := range results {
			if res.ConfigPath == globalPath {
				foundGlobal = true
				if res.Status != model.StatusPassed {
					t.Errorf("Global result should pass")
				}
			}
			if res.ConfigPath == localPath {
				foundLocal = true
				if res.Status != model.StatusFailed {
					t.Errorf("Local result should fail")
				}
			}
		}

		if !foundGlobal || !foundLocal {
			t.Errorf("Did not find both global and local results (Global: %v, Local: %v)", foundGlobal, foundLocal)
		}
	})
}
