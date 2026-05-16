package audit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shiro/dependency-shield/internal/config"
	"github.com/shiro/dependency-shield/internal/model"
)

func TestAuditNpm(t *testing.T) {
	tmpDir := t.TempDir()
	p := config.NewDefaultPolicy()
	
	// Test File Not Found
	res := AuditNpm(filepath.Join(tmpDir, "nonexistent"), p)
	if res.Status != model.StatusSkip {
		t.Errorf("Expected StatusSkip, got %s", res.Status)
	}

	// Test Passed
	pathPassed := filepath.Join(tmpDir, ".npmrc_passed")
	os.WriteFile(pathPassed, []byte("min-release-age=30\n"), 0644)
	res = AuditNpm(pathPassed, p)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed (wrong value)
	pathFailed := filepath.Join(tmpDir, ".npmrc_failed")
	os.WriteFile(pathFailed, []byte("min-release-age=10\n"), 0644)
	res = AuditNpm(pathFailed, p)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}

	// Test Failed (missing key)
	pathMissing := filepath.Join(tmpDir, ".npmrc_missing")
	os.WriteFile(pathMissing, []byte("other-key=30\n"), 0644)
	res = AuditNpm(pathMissing, p)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}

	// Test Custom Policy
	p7 := config.Policy{MinAgeDays: 7}
	path7 := filepath.Join(tmpDir, ".npmrc_7")
	os.WriteFile(path7, []byte("min-release-age=7\n"), 0644)
	res = AuditNpm(path7, p7)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed for 7 days, got %s", res.Status)
	}

	// Test High security configuration (30 days) against low security policy (7 days)
	path30 := filepath.Join(tmpDir, ".npmrc_30")
	os.WriteFile(path30, []byte("min-release-age=30\n"), 0644)
	res = AuditNpm(path30, p7)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed for 30 days config against 7 days policy, got %s", res.Status)
	}
}

func TestAuditPnpm(t *testing.T) {
	tmpDir := t.TempDir()
	p := config.NewDefaultPolicy()

	// Test Passed
	pathPassed := filepath.Join(tmpDir, ".pnpmrc_passed")
	os.WriteFile(pathPassed, []byte("minimum-release-age=43200\n"), 0644)
	res := AuditPnpm(pathPassed, p)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Passed (higher value)
	pathPassedHigher := filepath.Join(tmpDir, ".pnpmrc_higher")
	os.WriteFile(pathPassedHigher, []byte("minimum-release-age=50000\n"), 0644)
	res = AuditPnpm(pathPassedHigher, p)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed
	pathFailed := filepath.Join(tmpDir, ".pnpmrc_failed")
	os.WriteFile(pathFailed, []byte("minimum-release-age=100\n"), 0644)
	res = AuditPnpm(pathFailed, p)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}

	// Test Custom Policy
	p7 := config.Policy{MinAgeDays: 7}
	path7 := filepath.Join(tmpDir, ".pnpmrc_7")
	os.WriteFile(path7, []byte("minimum-release-age=10080\n"), 0644) // 7 * 24 * 60
	res = AuditPnpm(path7, p7)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed for 7 days (10080 mins), got %s", res.Status)
	}
}

func TestAuditUv(t *testing.T) {
	tmpDir := t.TempDir()
	p := config.NewDefaultPolicy()

	// Test Passed (tool.uv)
	pathPassed := filepath.Join(tmpDir, "uv.toml_passed")
	os.WriteFile(pathPassed, []byte("[tool.uv]\nexclude-newer = \"30d\"\n"), 0644)
	res := AuditUv(pathPassed, p)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Passed (top-level)
	pathPassedTop := filepath.Join(tmpDir, "uv.toml_passed_top")
	os.WriteFile(pathPassedTop, []byte("exclude-newer = \"30d\"\n"), 0644)
	res = AuditUv(pathPassedTop, p)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed
	pathFailed := filepath.Join(tmpDir, "uv.toml_failed")
	os.WriteFile(pathFailed, []byte("exclude-newer = \"7d\"\n"), 0644)
	res = AuditUv(pathFailed, p)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}

	// Test Custom Policy
	p7 := config.Policy{MinAgeDays: 7}
	path7 := filepath.Join(tmpDir, "uv.toml_7")
	os.WriteFile(path7, []byte("exclude-newer = \"7d\"\n"), 0644)
	res = AuditUv(path7, p7)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed for 7 days, got %s", res.Status)
	}

	// Test High security configuration (30 days) against low security policy (7 days)
	path30 := filepath.Join(tmpDir, "uv.toml_30")
	os.WriteFile(path30, []byte("exclude-newer = \"30d\"\n"), 0644)
	res = AuditUv(path30, p7)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed for 30 days config against 7 days policy, got %s", res.Status)
	}
}

func TestAuditBun(t *testing.T) {
	tmpDir := t.TempDir()
	p := config.NewDefaultPolicy()

	// Test Passed
	pathPassed := filepath.Join(tmpDir, "bunfig.toml_passed")
	os.WriteFile(pathPassed, []byte("[install]\nminimumReleaseAge = 2592000\n"), 0644)
	res := AuditBun(pathPassed, p)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Passed (higher value)
	pathPassedHigher := filepath.Join(tmpDir, "bunfig.toml_higher")
	os.WriteFile(pathPassedHigher, []byte("[install]\nminimumReleaseAge = 3000000\n"), 0644)
	res = AuditBun(pathPassedHigher, p)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed, got %s", res.Status)
	}

	// Test Failed
	pathFailed := filepath.Join(tmpDir, "bunfig.toml_failed")
	os.WriteFile(pathFailed, []byte("[install]\nminimumReleaseAge = 100\n"), 0644)
	res = AuditBun(pathFailed, p)
	if res.Status != model.StatusFailed {
		t.Errorf("Expected StatusFailed, got %s", res.Status)
	}

	// Test Custom Policy
	p7 := config.Policy{MinAgeDays: 7}
	path7 := filepath.Join(tmpDir, "bunfig.toml_7")
	os.WriteFile(path7, []byte("[install]\nminimumReleaseAge = 604800\n"), 0644) // 7 * 24 * 3600
	res = AuditBun(path7, p7)
	if res.Status != model.StatusPassed {
		t.Errorf("Expected StatusPassed for 7 days (604800 secs), got %s", res.Status)
	}
}

func TestAuditTool(t *testing.T) {
	tmpDir := t.TempDir()
	p := config.NewDefaultPolicy()

	path1 := filepath.Join(tmpDir, ".npmrc_1")
	os.WriteFile(path1, []byte("min-release-age=30\n"), 0644)
	path2 := filepath.Join(tmpDir, ".npmrc_2")
	os.WriteFile(path2, []byte("min-release-age=10\n"), 0644)

	results := AuditTool("npm", []string{path1, path2}, p)

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
