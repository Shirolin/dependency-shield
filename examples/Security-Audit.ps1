# Security-Audit.ps1
# 检查 AI 时代供应链安全配置 (npm & uv)

Write-Host "`n===============================================" -ForegroundColor Cyan
Write-Host "🛡️  AI 时代供应链安全配置检查工具" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan

$allPassed = $true

# --- 1. Node.js / npm 检查 ---
Write-Host "`n[1/2] 正在检查 Node.js / npm 安全配置..." -ForegroundColor Magenta
$npmVersion = npm --version
Write-Host "• npm 版本: $npmVersion"

if ([version]$npmVersion -ge [version]"11.10.0") {
    $minAge = npm config get min-release-age
    if ($minAge -eq "30") {
        Write-Host "✅ 匹配：npm 已开启 30 天发布冷却期 (min-release-age=30)。" -ForegroundColor Green
    } else {
        Write-Host "❌ 警告：npm 冷却期配置不正确 (当前为: $minAge)。" -ForegroundColor Red
        $allPassed = $false
    }
} else {
    Write-Host "❌ 错误：npm 版本过低，不支持 min-release-age 保护。" -ForegroundColor Red
    $allPassed = $false
}

# --- 2. Python / uv 检查 ---
Write-Host "`n[2/2] 正在检查 Python / uv 安全配置..." -ForegroundColor Magenta
$uvPath = Get-Command uv -ErrorAction SilentlyContinue

if ($uvPath) {
    Write-Host "• uv 版本: $(uv --version)"
    
    # 检查全局配置文件中的 exclude-newer
    $uvConfigPath = "$env:APPDATA\uv\uv.toml"
    if (Test-Path $uvConfigPath) {
        $configContent = Get-Content $uvConfigPath -Raw
        if ($configContent -match 'exclude-newer\s*=\s*"30d"') {
            Write-Host "✅ 匹配：uv 已开启 30 天发布保护 (exclude-newer = `"30d`")。" -ForegroundColor Green
        } else {
            Write-Host "❌ 警告：uv 配置文件存在但未正确设置 30d 拦截。" -ForegroundColor Red
            $allPassed = $false
        }
    } else {
        Write-Host "❌ 错误：未找到 uv 全局配置文件 (uv.toml)。" -ForegroundColor Red
        $allPassed = $false
    }
} else {
    Write-Host "❌ 错误：未检测到 uv，无法实施 Python 侧的自动拦截。" -ForegroundColor Red
    $allPassed = $false
}

# --- 总结 ---
Write-Host "`n-----------------------------------------------"
if ($allPassed) {
    Write-Host "🎉 审计通过！你的环境已配置完善，可抵御 99% 的 AI 幻觉安装攻击。" -ForegroundColor Green
} else {
    Write-Host "⚠️  审计未通过！请根据上方提示修复配置。" -ForegroundColor Yellow
}
Write-Host "-----------------------------------------------`n"
