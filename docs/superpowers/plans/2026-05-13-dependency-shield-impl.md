# DependencyShield Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个基于 Go 的零依赖 CLI 工具，用于审计和修复 npm, pnpm, uv, bun 的“发布冷却期”配置，以防御供应链攻击。

**Architecture:** 采用模块化设计，包括环境侦测 (Prober)、审计引擎 (Audit)、配置修复 (Fixer) 和 Cobra CLI 前端。优先使用正则表达式进行非侵入式的文件修改。

**Tech Stack:** Go 1.21+, Cobra (CLI), Fatih Color (UI), Go-toml/v2 (TOML parsing).

---

## File Structure
- `go.mod`: 项目依赖管理。
- `main.go`: 程序入口。
- `internal/config/constants.go`: 定义全局常量、路径和策略阈值。
- `internal/model/model.go`: 定义审计结果等核心数据结构。
- `internal/prober/prober.go`: 环境侦测逻辑（定位配置文件）。
- `internal/audit/engine.go`: 审计逻辑。
- `internal/fixer/fixer.go`: 配置文件修改逻辑。
- `cmd/root.go`, `cmd/audit.go`, `cmd/fix.go`: CLI 命令定义。

---

### Task 1: Project Initialization

**Files:**
- Create: `go.mod`
- Create: `main.go`

- [ ] **Step 1: Initialize Go module**
Run: `go mod init github.com/youruser/dependencyshield`

- [ ] **Step 2: Install dependencies**
Run: `go get github.com/spf13/cobra github.com/fatih/color github.com/pelletier/go-toml/v2`

- [ ] **Step 3: Create minimal main.go**
```go
package main

import "github.com/youruser/dependencyshield/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 4: Commit**
```bash
git add go.mod go.sum main.go
git commit -m "chore: initialize go project"
```

---

### Task 2: Define Core Models and Constants

**Files:**
- Create: `internal/config/constants.go`
- Create: `internal/model/model.go`

- [ ] **Step 1: Define Constants**
```go
package config

const (
	DefaultMinAgeDays = 30
	NpmMinAge         = "30"
	PnpmMinAgeMins    = "43200"    // 30 * 24 * 60
	BunMinAgeSecs     = "2592000"  // 30 * 24 * 3600
	UvExcludeNewer    = "30d"
)
```

- [ ] **Step 2: Define Models**
```go
package model

type ToolStatus string

const (
	StatusPassed ToolStatus = "PASSED"
	StatusFailed ToolStatus = "FAILED"
	StatusWarn   ToolStatus = "WARNING"
	StatusSkip   ToolStatus = "SKIPPED"
)

type AuditResult struct {
	ToolName    string
	ConfigPath  string
	CurrentVal  string
	Status      ToolStatus
	Message     string
}
```

- [ ] **Step 3: Commit**
```bash
git add internal/config/constants.go internal/model/model.go
git commit -m "feat: define core models and constants"
```

---

### Task 3: Implement Environment Prober

**Files:**
- Create: `internal/prober/prober.go`
- Test: `internal/prober/prober_test.go`

- [ ] **Step 1: Write Prober logic to locate configs**
```go
package prober

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetNpmrcPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".npmrc")
}

func GetUvConfigPath() string {
	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		return filepath.Join(appdata, "uv", "uv.toml")
	}
	// Simplified for brevity, add XDG logic in actual impl
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "uv", "uv.toml")
}
```

- [ ] **Step 2: Write tests for Prober**
```go
package prober

import "testing"

func TestPaths(t *testing.T) {
    p := GetNpmrcPath()
    if p == "" { t.Error("Path should not be empty") }
}
```

- [ ] **Step 3: Commit**
```bash
git add internal/prober/
git commit -m "feat: implement environment prober"
```

---

### Task 4: Implement Audit Engine

**Files:**
- Create: `internal/audit/engine.go`
- Test: `internal/audit/engine_test.go`

- [ ] **Step 1: Implement Npm Audit logic**
```go
package audit

import (
	"bufio"
	"os"
	"strings"
	"github.com/youruser/dependencyshield/internal/model"
)

func AuditNpm(path string) model.AuditResult {
	res := model.AuditResult{ToolName: "npm", ConfigPath: path, Status: model.StatusFailed}
	file, err := os.Open(path)
	if err != nil {
		res.Status = model.StatusSkip
		return res
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "min-release-age=30") {
			res.Status = model.StatusPassed
			return res
		}
	}
	return res
}
```

- [ ] **Step 2: Write tests and commit**
```bash
git add internal/audit/
git commit -m "feat: implement audit engine"
```

---

### Task 5: Implement Fixer Service (Regex-based)

**Files:**
- Create: `internal/fixer/fixer.go`

- [ ] **Step 1: Implement safe file writing with regex replacement**
```go
package fixer

import (
	"os"
	"regexp"
)

func FixNpmrc(path string) error {
	content, _ := os.ReadFile(path)
	re := regexp.MustCompile(`(?m)^min-release-age=.*$`)
	if re.Match(content) {
		newContent := re.ReplaceAll(content, []byte("min-release-age=30"))
		return os.WriteFile(path, newContent, 0644)
	}
	// Append if not found
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	_, err := f.WriteString("\nmin-release-age=30\n")
	return err
}
```

- [ ] **Step 2: Commit**
```bash
git add internal/fixer/
git commit -m "feat: implement fixer service"
```

---

### Task 6: Implement CLI Commands (Cobra)

**Files:**
- Create: `cmd/root.go`
- Create: `cmd/audit.go`
- Create: `cmd/fix.go`

- [ ] **Step 1: Setup Root command**
- [ ] **Step 2: Implement Audit command with output formatting**
- [ ] **Step 3: Implement Fix command with confirmation**
- [ ] **Step 4: Commit and Final Test**
```bash
go build -o shield
./shield audit
```
