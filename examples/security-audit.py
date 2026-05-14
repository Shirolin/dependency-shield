import os
import platform
import subprocess
import shutil
import re
from pathlib import Path

# --- 配置常量 ---
TARGET_DAYS = 30
PNPM_MINS = TARGET_DAYS * 24 * 60
BUN_SECS = TARGET_DAYS * 24 * 3600

def print_header(text):
    print(f"\n{'='*50}")
    print(f"🛡️  {text}")
    print(f"{'='*50}")

def get_npmrc_path():
    return Path.home() / ".npmrc"

def get_uv_config_path():
    if platform.system() == "Windows":
        # Windows 优先检查 APPDATA
        appdata = os.environ.get("APPDATA")
        if appdata:
            return Path(appdata) / "uv" / "uv.toml"
        return Path.home() / "AppData" / "Roaming" / "uv" / "uv.toml"
    else:
        # Linux/macOS 遵循 XDG 标准
        xdg_config = os.environ.get("XDG_CONFIG_HOME")
        if xdg_config:
            return Path(xdg_config) / "uv" / "uv.toml"
        return Path.home() / ".config" / "uv" / "uv.toml"

def check_npm():
    print("\n[1/4] Node.js / npm 检查...")
    npm_path = shutil.which("npm")
    if not npm_path:
        print("ℹ️ 未检测到 npm")
        return True

    try:
        version_cmd = ["npm", "--version"]
        if platform.system() == "Windows":
            version_cmd = ["cmd", "/c", "npm", "--version"]
        version = subprocess.check_output(version_cmd, text=True, stderr=subprocess.STDOUT).strip()
        print(f"• npm 版本: {version}")
        
        npmrc = get_npmrc_path()
        if npmrc.exists():
            content = npmrc.read_text(encoding="utf-8")
            if "min-release-age=30" in content:
                print("✅ 匹配：.npmrc 已开启 30 天发布冷却期。")
                return True
        print("❌ 警告：npm 未在 .npmrc 中发现正确的 30 天拦截配置。")
        return False
    except Exception as e:
        print(f"⚠️ 检查过程中出错: {e}")
        return False

def check_uv():
    print("\n[2/4] Python / uv 检查...")
    uv_path = shutil.which("uv")
    if not uv_path:
        print("ℹ️ 未检测到 uv")
        return True

    try:
        version = subprocess.check_output(["uv", "--version"], text=True).strip()
        print(f"• {version}")
        
        config_path = get_uv_config_path()
        if config_path.exists():
            content = config_path.read_text(encoding="utf-8")
            if re.search(r'exclude-newer\s*=\s*"30d"', content):
                print("✅ 匹配：uv 已开启 30 天发布保护。")
                return True
        print(f"❌ 警告：未在 {config_path} 发现 30d 拦截配置。")
        return False
    except Exception as e:
        print(f"⚠️ 检查失败: {e}")
        return False

def check_pnpm():
    print("\n[3/4] pnpm 检查...")
    pnpm_path = shutil.which("pnpm")
    if not pnpm_path:
        print("ℹ️ 未检测到 pnpm")
        return True

    try:
        cmd = ["pnpm", "config", "get", "minimum-release-age"]
        if platform.system() == "Windows":
            cmd = ["cmd", "/c", "pnpm", "config", "get", "minimum-release-age"]
        val = subprocess.check_output(cmd, text=True).strip()
        if val != "undefined" and val.isdigit() and int(val) >= PNPM_MINS:
            print(f"✅ 匹配：pnpm 已开启 {val} 分钟冷却期 (>=30天)。")
            return True
        print(f"❌ 警告：pnpm 冷却期不足 30 天 (当前: {val})。")
        return False
    except Exception as e:
        print(f"⚠️ 检查失败: {e}")
        return False

def check_bun():
    print("\n[4/4] Bun 检查...")
    bun_path = shutil.which("bun")
    if not bun_path:
        print("ℹ️ 未检测到 bun")
        return True

    bunfig = Path.home() / ".bunfig.toml"
    if bunfig.exists():
        content = bunfig.read_text(encoding="utf-8")
        if f"minimumReleaseAge = {BUN_SECS}" in content or "minimumReleaseAge = 2592000" in content:
            print("✅ 匹配：Bun 已开启 30 天发布冷却期。")
            return True
    print("❌ 警告：Bun 配置文件未设置正确的 30 天拦截。")
    return False

def main():
    print_header(f"跨平台供应链安全审计 (系统: {platform.system()})")
    
    results = [
        check_npm(),
        check_uv(),
        check_pnpm(),
        check_bun()
    ]

    print("\n" + "-"*50)
    if all(results):
        print("🎉 审计通过！你的环境具备跨工具的 AI 幻觉防御能力。")
    else:
        print("⚠️  审计未完全通过！请检查上方红色警告。")
    print("-"*50 + "\n")

if __name__ == "__main__":
    main()
