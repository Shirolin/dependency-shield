# 修正审计阈值判断逻辑计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将审计逻辑从“精确匹配”修改为“大于等于匹配”，确保更严格的安全设置（如 30 天）能通过较低要求（如 7 天）的审计。

**Architecture:** 
1. 修改 `internal/audit/engine.go` 中的 `AuditNpm` 和 `AuditUv` 函数。
2. 为 `AuditNpm` 添加数值解析，将读取到的字符串值转换为整数进行比较。
3. 为 `AuditUv` 添加对 `exclude-newer` 值的解析（解析如 "30d" 中的数字）。
4. 更新 `internal/audit/engine_test.go` 以验证这一逻辑。

**Tech Stack:** Go

---

### Task 1: 修正 AuditNpm 逻辑

**Files:**
- Modify: `internal/audit/engine.go`

- [ ] **Step 1: 修改 AuditNpm 以支持数值比较**

```go
// internal/audit/engine.go

// AuditNpm scans for 'min-release-age' in the given path using the provided policy.
func AuditNpm(path string, p config.Policy) model.AuditResult {
    // ... (文件打开逻辑保持不变)
    
	target := p.MinAgeDays // 直接使用天数整数比较
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "min-release-age") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				valStr := strings.TrimSpace(parts[1])
				result.CurrentVal = valStr
                val, err := strconv.Atoi(valStr)
				if err == nil && val >= target {
					result.Status = model.StatusPassed
					result.Message = "Security policy met"
					return result
				}
			}
		}
	}

	if result.CurrentVal == "" {
		result.Message = "min-release-age not found"
	} else {
		result.Message = "min-release-age is less than " + strconv.Itoa(target)
	}

	return result
}
```

- [ ] **Step 2: 提交更改**

```bash
git add internal/audit/engine.go
git commit -m "fix(audit): change npm threshold check to greater-than-or-equal"
```

### Task 2: 修正 AuditUv 逻辑

**Files:**
- Modify: `internal/audit/engine.go`

- [ ] **Step 1: 修改 AuditUv 解析逻辑**

```go
// internal/audit/engine.go

// 需要解析 "30d" 这种格式
func AuditUv(path string, p config.Policy) model.AuditResult {
    // ... (TOML 解析逻辑保持不变)

	target := p.MinAgeDays
	// ... (获取 val 的逻辑保持不变)

	if ok {
		if s, ok := val.(string); ok {
			result.CurrentVal = s
            // 解析数字部分，例如 "30d" -> 30
            numStr := strings.TrimSuffix(s, "d")
            valNum, err := strconv.Atoi(numStr)
			if err == nil && valNum >= target {
				result.Status = model.StatusPassed
				result.Message = "Security policy met"
				return result
			}
		}
	}

	if !ok {
		result.Message = "exclude-newer not found"
	} else {
		result.Message = "exclude-newer is less than " + strconv.Itoa(target) + "d"
	}

	return result
}
```

- [ ] **Step 2: 提交更改**

```bash
git add internal/audit/engine.go
git commit -m "fix(audit): change uv threshold check to greater-than-or-equal"
```

### Task 3: 更新单元测试并验证

**Files:**
- Modify: `internal/audit/engine_test.go`
- Test: `go test ./internal/audit/...`

- [ ] **Step 1: 添加测试用例：30天配置应通过 7天审计**
- [ ] **Step 2: 运行测试并验证**

Run: `go test ./internal/audit/... -v`
Expected: PASS
