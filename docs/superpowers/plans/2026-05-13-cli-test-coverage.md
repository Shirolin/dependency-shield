# DependencyShield CLI Test Coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 补齐 `cmd` 包的测试覆盖率，确保 CLI 命令及其参数（如 `--force`）工作正常。

**Architecture:** 
1. **重构解耦**: 将 `cmd` 中的输出逻辑从全局 `os.Stdout` 重构为支持注入 `io.Writer`。
2. **命令测试**: 编写针对 `auditCmd` 和 `fixCmd` 的单元测试，捕获并验证其输出字符串。
3. **参数验证**: 测试 `fix` 命令在不同 flag 下的行为。

**Tech Stack:** Go 1.21+, testing (standard library), bytes (for buffer).

---

### Task 1: Refactor CLI for Testability

**Files:**
- Modify: `cmd/root.go`
- Modify: `cmd/audit.go`
- Modify: `cmd/fix.go`

- [ ] **Step 1: Update RootCmd to support output injection**
在 `cmd` 包中增加一个可配置的输出源。

```go
var outWriter io.Writer = os.Stdout

func SetOut(w io.Writer) {
    outWriter = w
}
```

- [ ] **Step 2: Update audit and fix commands to use outWriter**
将 `fmt.Printf` 和 `color.Printf` 替换为支持 `io.Writer` 的调用。

- [ ] **Step 3: Commit**
```bash
git commit -m "refactor: support output injection in cmd package for testing"
```

---

### Task 2: Implement Audit Command Tests

**Files:**
- Create: `cmd/audit_test.go`

- [ ] **Step 1: Write a failing test for audit command**
```go
func TestAuditCommand(t *testing.T) {
    buf := new(bytes.Buffer)
    SetOut(buf)
    // 执行 auditCmd.Run()
    // 断言 buf.String() 包含 "🛡️ DependencyShield Audit Report"
}
```

- [ ] **Step 2: Implement and Verify GREEN**

---

### Task 3: Implement Fix Command & Flag Tests

**Files:**
- Create: `cmd/fix_test.go`

- [ ] **Step 1: Write tests for fix command**
验证：
1. 正常运行输出。
2. `--force` flag 被正确解析（虽然目前逻辑中未强制中断，但需验证 flag 存在）。

- [ ] **Step 2: Implement and Verify GREEN**

---

### Task 4: Final Coverage Verification

- [ ] **Step 1: Run coverage again**
```bash
go test ./cmd/... -cover
```
预期：`cmd` 包覆盖率从 0% 提升至 70% 以上。

- [ ] **Step 2: Commit & Done**
```bash
git commit -m "test: achieve high coverage for cmd package"
```
