package internal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shiro/dependencyshield/internal/audit"
	"github.com/shiro/dependencyshield/internal/model"
	"github.com/shiro/dependencyshield/internal/prober"
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
	// 5. Currently, our 'audit' and 'prober' packages only support global paths. Therefore, this test SHOULD FAIL.
	t.Run("Automated Discovery", func(t *testing.T) {
		// In a real application run, we would trigger an audit that should find both.
		// Currently, we only have GetNpmrcPath which returns a single string (the global one).
		
		results := []model.AuditResult{
			audit.AuditNpm(prober.GetNpmrcPath()),
		}

		// We expect 2 results if hierarchy is supported, but we only get 1 (the global one).
		if len(results) < 2 {
			t.Errorf("Expected at least 2 audit results for npm (global and local), but found %d. Configuration hierarchy support is missing.", len(results))
		}
	})
}
