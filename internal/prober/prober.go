package prober

import (
	"os"
	"path/filepath"
	"runtime"
)

// FindConfigUpwards searches for a file with the given name starting from the current
// working directory and moving upwards to the root directory.
func FindConfigUpwards(fileName string) string {
	curr, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		path := filepath.Join(curr, fileName)
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	return ""
}

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

// GetLocalNpmrcPath returns the path to the local '.npmrc' file if found in the directory hierarchy.
func GetLocalNpmrcPath() string {
	return FindConfigUpwards(".npmrc")
}

// GetLocalPnpmrcPath returns the path to the local '.npmrc' file (used by pnpm) if found.
func GetLocalPnpmrcPath() string {
	return FindConfigUpwards(".npmrc")
}

// GetLocalUvConfigPath returns the path to the local 'uv.toml' file if found.
func GetLocalUvConfigPath() string {
	return FindConfigUpwards("uv.toml")
}

// GetLocalBunfigPath returns the path to the local '.bunfig.toml' file if found.
func GetLocalBunfigPath() string {
	return FindConfigUpwards(".bunfig.toml")
}
