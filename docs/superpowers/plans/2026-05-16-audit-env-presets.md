# 环境预设与建议功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为审计功能添加环境预设（Local, CI, Prod）支持，并根据审计结果提供安全建议，引导用户选择合适的冷却期。

**Architecture:** 
1. 在 `internal/config` 中定义环境预设常量及其对应的冷却期天数。
2. 扩展 `cmd/audit.go`，添加 `--env` (或 `-e`) 标志，允许用户通过环境名称选择策略。
3. 修改审计报告输出逻辑，在报告底部添加基于当前结果和环境建议的提示信息。

**Tech Stack:** Go, Cobra (CLI framework)

---

### Task 1: 定义环境预设

**Files:**
- Modify: `internal/config/constants.go`

- [ ] **Step 1: 在 constants.go 中添加环境预设定义**

```go
package config

import (
	"fmt"
	"strconv"
)

const (
	DefaultMinAgeDays = 30
	// 环境预设天数
	EnvLocalDays = 7
	EnvCIDays    = 15
	EnvProdDays  = 30
)

// Environment presets
const (
	EnvLocal = "local"
	EnvCI    = "ci"
	EnvProd  = "prod"
)

// ... (rest of the file)
```

- [ ] **Step 2: 提交更改**

```bash
git add internal/config/constants.go
git commit -m "feat(config): add environment presets for min-release-age"
```

### Task 2: 扩展 Audit 命令支持环境选择

**Files:**
- Modify: `cmd/audit.go`

- [ ] **Step 1: 添加 env 变量并更新 init 函数**

```go
var env string

func init() {
	auditCmd.Flags().IntVarP(&minAgeDays, "days", "d", config.DefaultMinAgeDays, "Minimum release age in days")
	auditCmd.Flags().StringVarP(&env, "env", "e", "", "Environment preset: local (7d), ci (15d), prod (30d). Overrides --days.")
	rootCmd.AddCommand(auditCmd)
}
```

- [ ] **Step 2: 在 Run 逻辑中处理环境预设覆盖**

```go
// 在 auditCmd 的 Run 函数开始处：
Run: func(cmd *cobra.Command, args []string) {
    // 根据环境预设覆盖 minAgeDays
    switch env {
    case config.EnvLocal:
        minAgeDays = config.EnvLocalDays
    case config.EnvCI:
        minAgeDays = config.EnvCIDays
    case config.EnvProd:
        minAgeDays = config.EnvProdDays
    case "":
        // 保持原样 (默认或通过 --days 指定)
    default:
        fmt.Fprintf(outWriter, "⚠️  Unknown environment '%s', falling back to %d days\n", env, minAgeDays)
    }

    fmt.Fprintln(outWriter, "🛡️  DependencyShield Audit Report")
    // ...
```

- [ ] **Step 3: 提交更改**

```bash
git add cmd/audit.go
git commit -m "feat(cmd): add --env flag to audit command"
```

### Task 3: 添加审计报告建议提示

**Files:**
- Modify: `cmd/audit.go`

- [ ] **Step 1: 在审计完成后输出建议信息**

```go
// 在 auditCmd 的 Run 函数末尾循环结束后：
		for _, t := range tools {
            // ... (现有的审计循环)
        }

        // 输出建议
        fmt.Fprintln(outWriter, "\n💡 Security Recommendations:")
        if minAgeDays < 15 {
            fmt.Fprintln(outWriter, "- [Local] Your policy is set to Local development (7-14 days). This balances speed and safety.")
        }
        if minAgeDays >= 15 && minAgeDays < 30 {
            fmt.Fprintln(outWriter, "- [CI/Test] Your policy is set to CI/Test (15 days). Recommended for shared staging environments.")
        }
        if minAgeDays >= 30 {
            fmt.Fprintln(outWriter, "- [Prod] Your policy is set to Production (30 days). Maximum protection against supply chain attacks.")
        }
        
        fmt.Fprintln(outWriter, "- Suggestion: Use '--env' to quickly switch presets: 'audit -e local', 'audit -e ci', or 'audit -e prod'.")
	},
```

- [ ] **Step 2: 提交更改**

```bash
git add cmd/audit.go
git commit -m "feat(cmd): add security recommendations to audit report"
```

### Task 4: 验证测试

- [ ] **Step 1: 验证默认行为 (30天)**

Run: `go run main.go audit`
Expected: 策略显示为 30 天，底部有 Prod 建议。

- [ ] **Step 2: 验证环境预设 (local/7天)**

Run: `go run main.go audit -e local`
Expected: 策略显示为 7 天，底部有 Local 建议，之前不满足 30 天但满足 7 天的项应显示为 PASSED (例如 demo_project 中的 npm)。

- [ ] **Step 3: 验证环境预设覆盖 --days**

Run: `go run main.go audit -e ci -d 30`
Expected: 策略显示为 15 天 (环境优先)。

- [ ] **Step 4: 运行所有测试确保无回归**

Run: `go test ./...`
Expected: PASS
