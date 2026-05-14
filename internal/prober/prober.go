package prober

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetNpmrcPath returns the path to the global '.npmrc' file.
func GetNpmrcPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".npmrc")
}

// GetPnpmrcPath returns the path to the global '.pnpmrc' file.
func GetPnpmrcPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".npmrc") // Instructions say "Same as npmrc", but usually it's .npmrc for pnpm too, or .pnpmrc. Let's stick to instructions.
}

// GetUvConfigPath returns the path to the global 'uv.toml' file.
func GetUvConfigPath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			// Fallback to home if APPDATA is not set
			home, err := os.UserHomeDir()
			if err != nil {
				return ""
			}
			return filepath.Join(home, "AppData", "Roaming", "uv", "uv.toml")
		}
		return filepath.Join(appData, "uv", "uv.toml")
	}

	// Linux/macOS
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, "uv", "uv.toml")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "uv", "uv.toml")
}

// GetBunfigPath returns the path to the global '.bunfig.toml' file.
func GetBunfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".bunfig.toml")
}
