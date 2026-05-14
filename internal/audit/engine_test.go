package audit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shiro/dependencyshield/internal/model"
)

func TestAuditNpm(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test File Not Found
	res := AuditNpm(filepath.Join(tmpDir, "nonexistent"))
	if res.Status != model.StatusSkip {
		t.Errorf("Expected StatusSkip, got %s", res.Status)
	}

	// Test Passed
	pathPassed := filepath.Join(tmpDir, ".npmrc_passed")
	os.WriteFile(pathPassed, []byte("min-release-age=30\n"), 0644)
	res = AuditNpm(pathPassed)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed (wrong value)
	pathFailed := filepath.Join(tmpDir, ".npmrc_failed")
	os.WriteFile(pathFailed, []byte("min-release-age=10\n"), 0644)
	res = AuditNpm(pathFailed)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}

	// Test Failed (missing key)
	pathMissing := filepath.Join(tmpDir, ".npmrc_missing")
	os.WriteFile(pathMissing, []byte("other-key=30\n"), 0644)
	res = AuditNpm(pathMissing)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}
}

func TestAuditPnpm(t *testing.T) {
	tmpDir := t.TempDir()

	// Test Passed
	pathPassed := filepath.Join(tmpDir, ".pnpmrc_passed")
	os.WriteFile(pathPassed, []byte("minimum-release-age=43200\n"), 0644)
	res := AuditPnpm(pathPassed)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Passed (higher value)
	pathPassedHigher := filepath.Join(tmpDir, ".pnpmrc_higher")
	os.WriteFile(pathPassedHigher, []byte("minimum-release-age=50000\n"), 0644)
	res = AuditPnpm(pathPassedHigher)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed
	pathFailed := filepath.Join(tmpDir, ".pnpmrc_failed")
	os.WriteFile(pathFailed, []byte("minimum-release-age=100\n"), 0644)
	res = AuditPnpm(pathFailed)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}
}

func TestAuditUv(t *testing.T) {
	tmpDir := t.TempDir()

	// Test Passed (tool.uv)
	pathPassed := filepath.Join(tmpDir, "uv.toml_passed")
	os.WriteFile(pathPassed, []byte("[tool.uv]\nexclude-newer = \"30d\"\n"), 0644)
	res := AuditUv(pathPassed)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Passed (top-level)
	pathPassedTop := filepath.Join(tmpDir, "uv.toml_passed_top")
	os.WriteFile(pathPassedTop, []byte("exclude-newer = \"30d\"\n"), 0644)
	res = AuditUv(pathPassedTop)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed
	pathFailed := filepath.Join(tmpDir, "uv.toml_failed")
	os.WriteFile(pathFailed, []byte("exclude-newer = \"7d\"\n"), 0644)
	res = AuditUv(pathFailed)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}
}

func TestAuditBun(t *testing.T) {
	tmpDir := t.TempDir()

	// Test Passed
	pathPassed := filepath.Join(tmpDir, "bunfig.toml_passed")
	os.WriteFile(pathPassed, []byte("[install]\nminimumReleaseAge = 2592000\n"), 0644)
	res := AuditBun(pathPassed)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Passed (higher value)
	pathPassedHigher := filepath.Join(tmpDir, "bunfig.toml_higher")
	os.WriteFile(pathPassedHigher, []byte("[install]\nminimumReleaseAge = 3000000\n"), 0644)
	res = AuditBun(pathPassedHigher)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed
	pathFailed := filepath.Join(tmpDir, "bunfig.toml_failed")
	os.WriteFile(pathFailed, []byte("[install]\nminimumReleaseAge = 100\n"), 0644)
	res = AuditBun(pathFailed)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}
}

func TestAuditTool(t *testing.T) {
	tmpDir := t.TempDir()

	path1 := filepath.Join(tmpDir, ".npmrc_1")
	os.WriteFile(path1, []byte("min-release-age=30\n"), 0644)
	path2 := filepath.Join(tmpDir, ".npmrc_2")
	os.WriteFile(path2, []byte("min-release-age=10\n"), 0644)

	results := AuditTool("npm", []string{path1, path2})

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].Status != model.StatusPassed {
		t.Errorf("Result 0: Expected StatusPassed, got %s", results[0].Status)
	}
	if results[1].Status != model.StatusFailed {
		t.Errorf("Result 1: Expected StatusFailed, got %s", results[1].Status)
	}
}
