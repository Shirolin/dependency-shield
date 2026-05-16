# Normalize Project Name Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Normalize the project name and module path to 'dependency-shield'.

**Architecture:** Update the module name in `go.mod`, replace all import paths in `.go` files, and update documentation links in `README.md`.

**Tech Stack:** Go (1.26.3)

---

### Task 1: Update go.mod

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Update module name**

```go
// From:
module github.com/shiro/dependencyshield
// To:
module github.com/shiro/dependency-shield
```

- [ ] **Step 2: Commit**

```bash
git add go.mod
git commit -m "refactor(config): update module name to dependency-shield in go.mod"
```

---

### Task 2: Update Go Imports

**Files:**
- Modify: `cmd/audit.go`
- Modify: `cmd/fix.go`
- Modify: `internal/audit/engine.go`
- Modify: `internal/audit/engine_test.go`
- Modify: `internal/fixer/fixer.go`
- Modify: `internal/integration_test.go`
- Modify: `main.go`

- [ ] **Step 1: Replace all occurrences of old module path**

Replace `github.com/shiro/dependencyshield` with `github.com/shiro/dependency-shield`.

- [ ] **Step 2: Commit**

```bash
git add cmd/ internal/ main.go
git commit -m "refactor: update import paths to dependency-shield"
```

---

### Task 3: Update README.md

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update GitHub links**

```markdown
// From:
[Releases](https://github.com/youruser/dependencyshield/releases)
// To:
[Releases](https://github.com/shiro/dependency-shield/releases)
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: update repository links in README.md"
```

---

### Task 4: Verify and Test

- [ ] **Step 1: Run all tests**

Run: `C:\Users\shiro\.vfox\cache\golang\v-1.26.3\golang-1.26.3\bin\go.exe test ./... -v`
Expected: ALL PASS

- [ ] **Step 2: Build the project**

Run: `C:\Users\shiro\.vfox\cache\golang\v-1.26.3\golang-1.26.3\bin\go.exe build -o shield.exe`
Expected: Success

- [ ] **Step 3: Commit (if any fixes were needed)**

```bash
git commit -m "chore: final adjustments after normalization"
```
