# DependencyShield 🛡️

**AI 时代下的供应链安全守护者**

DependencyShield 是一个专为个人开发者设计的轻量级、零依赖安全工具。它旨在防御恶意包安装和 AI 幻觉引导的依赖注入，通过审计和自动化配置主流包管理器的“发布冷却期”（Release Age Protection）功能来保护您的开发环境。

## ✨ 核心特性

- **零依赖运行**：单一静态二进制文件，无需预装运行时。
- **跨平台支持**：支持 Windows, macOS, Linux。
- **多工具覆盖**：支持 `npm`, `pnpm`, `uv`, `bun`。
- **项目级审计**：不仅检查全局配置，还支持向上递归审计项目级配置。
- **非侵入式修复**：使用正则技术保留您原有的配置注释和格式。

## 🚀 快速开始

### 安装
在 [Releases](https://github.com/shiro/dependency-shield/releases) 页面下载对应系统的二进制文件。

### 审计 (Audit)
检查您的开发环境是否安全：
```bash
./shield audit
```

### 修复 (Fix)
自动加固不合规的配置：
```bash
./shield fix
```

## 🛡️ 防御逻辑

DependencyShield 会确保以下配置被正确设置（默认 30 天）：
- **npm**: `min-release-age=30`
- **pnpm**: `minimum-release-age=43200` (分钟)
- **uv**: `exclude-newer = "30d"`
- **bun**: `minimumReleaseAge = 2592000` (秒)

## 🛠️ 构建与测试

如果您希望从源码构建：
```bash
go build -o shield.exe
go test ./... -v
```

## 📜 许可证

MIT License
