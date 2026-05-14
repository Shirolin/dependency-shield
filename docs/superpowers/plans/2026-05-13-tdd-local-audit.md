# DependencyShield TDD & Local Audit Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 引入 TDD 开发模式，实现项目级（Local）配置审计功能，并建立自动化回归测试套件。

**Architecture:** 
1. **测试先行**: 每个功能点先写集成测试/单元测试。
2. **递归搜索**: `prober` 模块增加向上递归搜索 `.npmrc` 等文件的逻辑。
3. **聚合报告**: CLI 汇总显示“全局状态”与“项目状态”。

**Tech Stack:** Go 1.21+, testing (standard library).

---

### Task 1: Setup Integration Test Environment

**Files:**
- Create: `internal/integration_test.go`

- [ ] **Step 1: Write a failing integration test**
创建一个模拟的文件系统结构：
- `tmp/mock_home/.npmrc` (配置正确)
- `tmp/mock_project/.npmrc` (配置错误)
测试 `shield audit` 是否能同时识别这两个层级。

```go
func TestAuditHierarchy(t *testing.T) {
    // 1. 设置模拟环境环境变量 (HOME, APPDATA等)
    // 2. 创建模拟文件
    // 3. 调用审计逻辑
    // 4. 断言：全局通过，局部失败
}
```

- [ ] **Step 2: Run test and verify RED**
Run: `go test ./internal/...`
Expected: FAIL (因为目前只查全局)

---

### Task 2: Implement Recursive Local Prober

**Files:**
- Modify: `internal/prober/prober.go`
- Test: `internal/prober/prober_test.go`

- [ ] **Step 1: Write failing test for local discovery**
```go
func TestFindLocalNpmrc(t *testing.T) {
    // 创建深层目录结构，在中间层放置 .npmrc
    // 断言 Prober 能从底层向上找到它
}
```

- [ ] **Step 2: Implement recursive search**
```go
func FindConfigUpwards(fileName string) string {
    curr, _ := os.Getwd()
    for {
        path := filepath.Join(curr, fileName)
        if _, err := os.Stat(path); err == nil { return path }
        parent := filepath.Dir(curr)
        if parent == curr { break }
        curr = parent
    }
    return ""
}
```

- [ ] **Step 3: Run tests and verify GREEN**

---

### Task 3: Update Audit Engine to handle slices of results

**Files:**
- Modify: `internal/model/model.go`
- Modify: `internal/audit/engine.go`

- [ ] **Step 1: Update AuditResult model to support multiple instances**
- [ ] **Step 2: Write tests for multi-file auditing**
- [ ] **Step 3: Implement logic to audit all found config files**

---

### Task 4: CLI Report Refactoring (Audit & Fix)

**Files:**
- Modify: `cmd/audit.go`
- Modify: `cmd/fix.go`

- [ ] **Step 1: Update UI to show "Global" vs "Local" sections**
- [ ] **Step 2: Update Fix command to fix ALL discovered files**
- [ ] **Step 3: Final regression test check**
运行所有测试，确保 `shield audit` 依然能正确检测最初的全局配置。

---
### Task 5: Commit & Final Verification

- [ ] **Step 1: Run all tests one last time**
- [ ] **Step 2: Final build and commit**
```bash
git commit -m "feat: add project-level auditing with TDD and regression tests"
```
