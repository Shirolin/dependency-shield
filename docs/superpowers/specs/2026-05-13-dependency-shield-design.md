# DependencyShield 设计说明书 (v1.0)

## 1. 概述 (Overview)
DependencyShield 是一个专为个人开发者设计的轻量级、零依赖安全工具，旨在防御 AI 时代下的供应链攻击（如恶意包安装、AI 幻觉引导的依赖注入）。它通过审计和自动化配置各大主流包管理器的“发布冷却期”（Release Age Protection）功能，确保开发环境不会在无意中安装发布时间极短（可能尚未被安全社区检测）的依赖包。

## 2. 目标 (Goals)
- **零依赖运行**：通过 Go 语言编译成单一静态二进制文件，用户无需预装 Python、Node.js 或任何运行库。
- **跨平台一致性**：无差别支持 Windows, macOS, 和主流 Linux 发行版。
- **安全与非侵入性**：在修改配置文件时，优先采用正则替换技术，最大限度保留用户的原有配置和注释。
- **极速体验**：秒级完成全系统环境扫描。

## 3. 架构设计 (Architecture)

### 3.1 核心组件
- **Env Prober (环境侦测器)**: 识别系统中安装的 npm, uv, pnpm, bun 等工具及其对应的配置文件位置。
- **Audit Engine (审计引擎)**: 读取配置文件并解析相关安全参数（如 `min-release-age`）。
- **Policy Manager (策略管理器)**: 维护“安全阈值”（默认 30 天/43200 分钟/2592000 秒）。
- **Fixer Service (修复服务)**: 交互式地纠正不合规的配置。
- **CLI Frontend**: 基于 `cobra` 实现的命令行界面。

### 3.2 技术栈
- **语言**: Go 1.21+
- **CLI 框架**: `github.com/spf13/cobra`
- **着色库**: `github.com/fatih/color`
- **解析库**: 针对 TOML (`uv.toml`, `bunfig.toml`) 使用 `github.com/pelletier/go-toml/v2`，针对 `.npmrc` 使用自定义正则引擎（以保留注释）。

## 4. 详细功能说明

### 4.1 审计逻辑 (Audit)
工具将扫描以下全局配置文件：
- **npm/pnpm**: `~/.npmrc` (Win: `%USERPROFILE%\.npmrc`)
  - 检查: `min-release-age` (npm) 和 `minimum-release-age` (pnpm, 单位: 分钟)。
- **uv**: `~/.config/uv/uv.toml` (Win: `%APPDATA%\uv\uv.toml`)
  - 检查: `exclude-newer = "30d"`。
- **Bun**: `~/.bunfig.toml` (Win: `%USERPROFILE%\.bunfig.toml`)
  - 检查: `[install] -> minimumReleaseAge = 2592000` (单位: 秒)。

### 4.2 修复逻辑 (Fix)
- 采用 **"Regex-First"** 策略：如果文件中已存在对应键名，使用正则表达式精准替换其数值，从而不影响同一行或相邻行的注释。
- 如果键名不存在，则在文件末尾（或相应 section 下）追加配置。

## 5. 交互界面 (Interface)

### 命令行指令
- `shield audit`: 扫描所有受支持的工具，生成合规报告。
- `shield fix`: 扫描并针对不合规项提供交互式修复选项（Y/N）。
- `shield version`: 显示版本信息。

### 报告示例
```text
🛡️ DependencyShield Audit Report
---------------------------------
[✅] npm: min-release-age=30
[❌] pnpm: minimum-release-age is NOT set
[✅] uv: exclude-newer="30d"
[⚠️] Bun: minimumReleaseAge=3600 (Current: 1h, Recommended: 30d)

Run 'shield fix' to resolve security issues.
```

## 6. 未来展望 (Roadmap)
- **v1.1**: 支持项目级（Local）配置文件的递归搜索与审计。
- **v1.2**: 增加针对特定包的“预安装检查”指令。
- **v2.0**: 扩展至其他安全领域（如镜像源安全、敏感信息扫描）。
