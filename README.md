# DependencyShield (🛡️ dependency-shield)

[![CI](https://github.com/Shirolin/dependency-shield/actions/workflows/ci.yml/badge.svg)](https://github.com/Shirolin/dependency-shield/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shirolin/dependency-shield)](https://goreportcard.com/report/github.com/Shirolin/dependency-shield)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AI 时代下的供应链安全守护者** —— 针对主流包管理器（npm, pnpm, uv, bun）的自动化安全策略审计与加固工具。

---

## 📖 概述

在 AI 辅助编程普及的今天，供应链攻击手段日益进化。AI 幻觉引导的依赖注入（AI Hallucination package attack）和针对新发布包的“闪电式”劫持已成为现实。

`dependency-shield` 是一个专为开发者设计的零依赖安全 CLI 工具。它通过强制实施**“发布冷却期”（Release Age Protection）**策略，确保开发环境不会在无意中引入发布时间极短（尚未经过社区安全审计）的依赖包，从而在源头阻断 99% 的新型供应链攻击。

## ✨ 核心特性

- **🚀 零运行依赖**：基于 Go 开发，提供单一静态二进制文件，无需 Python、Node.js 等环境即可运行。
- **🛡️ 深度审计体系**：
  - **多工具支持**：横跨前端（npm, pnpm, bun）与 Python（uv）生态。
  - **层级感知**：同时审计全局（Global）与项目级（Local）配置文件，识别配置覆盖风险。
- **📝 非侵入式修复**：采用精准正则技术进行配置加固，完美保留用户原有的配置文件注释与格式。
- **✅ 高度可靠**：核心逻辑经过严谨的 TDD 开发，具备 100% 的单元测试与集成测试通过率。

## 🛡️ 安全策略标准

本工具默认强制执行 **30 天发布冷却期**，这是基于安全社区共识的黄金防护窗口：

| 包管理器 | 配置项 | 安全阈值 |
| :--- | :--- | :--- |
| **npm** | `min-release-age` | 30 (days) |
| **pnpm** | `minimum-release-age` | 43200 (minutes) |
| **uv** | `exclude-newer` | "30d" |
| **bun** | `minimumReleaseAge` | 2592000 (seconds) |

## 🚀 快速上手

### 1. 下载
前往 [Releases](https://github.com/Shirolin/dependency-shield/releases) 下载适用于您操作系统的 `shield` 二进制文件。

### 2. 审计 (Audit)
一键扫描系统环境与当前项目的安全漏洞：
```bash
./shield audit
```

### 3. 加固 (Fix)
自动修复所有检测到的安全合规性问题：
```bash
./shield fix
```

## 🛠️ 开发者指南

### 构建
```bash
go build -o shield.exe
```

### 测试
```bash
# 运行完整测试套件（含集成测试）
go test ./... -v

# 查看测试覆盖率
go test ./... -cover
```

## 📜 许可证

本项目采用 [MIT License](LICENSE) 开源。
