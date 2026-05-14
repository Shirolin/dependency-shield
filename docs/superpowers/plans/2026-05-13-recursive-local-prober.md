# Recursive Local Prober Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement recursive upward search for local configuration files (.npmrc, .pnpmrc, uv.toml, .bunfig.toml).

**Architecture:** Add a generic `FindConfigUpwards` helper that traverses from CWD to root, then specialized wrappers for each config type.

**Tech Stack:** Go (Standard Library)

---

### Task 1: Test Recursive Upward Search

**Files:**
- Modify: `internal/prober/prober_test.go`

- [ ] **Step 1: Write the failing test `TestFindConfigUpwards`**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/prober -v -run TestFindConfigUpwards`
Expected: FAIL (undefined: FindConfigUpwards)

### Task 2: Implement FindConfigUpwards

**Files:**
- Modify: `internal/prober/prober.go`

- [ ] **Step 1: Implement `FindConfigUpwards`**

```go
// FindConfigUpwards searches for a file with the given name starting from the current
// working directory and moving upwards to the root directory.
func FindConfigUpwards(fileName string) string {
	curr, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		path := filepath.Join(curr, fileName)
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	return ""
}
```

- [ ] **Step 2: Add specialized GetLocal... functions**

```go
// GetLocalNpmrcPath returns the path to the local '.npmrc' file if found in the directory hierarchy.
func GetLocalNpmrcPath() string {
	return FindConfigUpwards(".npmrc")
}

// GetLocalPnpmrcPath returns the path to the local '.npmrc' file (used by pnpm) if found.
func GetLocalPnpmrcPath() string {
	return FindConfigUpwards(".npmrc")
}

// GetLocalUvConfigPath returns the path to the local 'uv.toml' file if found.
func GetLocalUvConfigPath() string {
	return FindConfigUpwards("uv.toml")
}

// GetLocalBunfigPath returns the path to the local '.bunfig.toml' file if found.
func GetLocalBunfigPath() string {
	return FindConfigUpwards(".bunfig.toml")
}
```

- [ ] **Step 3: Run test to verify it passes**

Run: `go test ./internal/prober -v`
Expected: PASS

### Task 3: Commit and Finish

- [ ] **Step 1: Commit changes**

```bash
git add internal/prober/prober.go internal/prober/prober_test.go
git commit -m "feat(prober): implement recursive upward search for local configs"
```
