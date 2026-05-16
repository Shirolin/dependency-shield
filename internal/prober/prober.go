package prober

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// IsToolInstalled checks if the given tool binary exists in the system PATH.
func IsToolInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

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

// FindConfigsDownwardsInDir searches for all files with the given name starting from the given
// root directory and moving downwards into subdirectories.
func FindConfigsDownwardsInDir(root, fileName string) []string {
	var paths []string
	
	ignoreDirs := map[string]bool{
		".git":         true,
		"node_modules": true,
		"venv":         true,
		".venv":        true,
		"target":       true,
		"dist":         true,
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if ignoreDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == fileName {
			absPath, _ := filepath.Abs(path)
			paths = append(paths, absPath)
		}
		return nil
	})

	if err != nil {
		return paths
	}

	return paths
}

// GetConfigsInPath returns all config files of a specific tool found within the given path (file or directory).
func GetConfigsInPath(toolName, path string) []string {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	fileName := ""
	switch toolName {
	case "npm", "pnpm":
		fileName = ".npmrc"
	case "uv":
		fileName = "uv.toml"
	case "bun":
		fileName = ".bunfig.toml"
	}

	if !info.IsDir() {
		if filepath.Base(path) == fileName {
			abs, _ := filepath.Abs(path)
			return []string{abs}
		}
		return nil
	}

	return FindConfigsDownwardsInDir(path, fileName)
}

// FindConfigsDownwards searches for all files with the given name starting from the current
// working directory and moving downwards into subdirectories.
func FindConfigsDownwards(fileName string) []string {
	root, err := os.Getwd()
	if err != nil {
		return nil
	}
	return FindConfigsDownwardsInDir(root, fileName)
}

// GetLocalNpmrcPaths returns all local '.npmrc' files found in the project.
func GetLocalNpmrcPaths() []string {
	return FindConfigsDownwards(".npmrc")
}

// GetLocalPnpmrcPaths returns all local '.npmrc' files (used by pnpm) found.
func GetLocalPnpmrcPaths() []string {
	return FindConfigsDownwards(".npmrc")
}

// GetLocalUvConfigPaths returns all local 'uv.toml' files found.
func GetLocalUvConfigPaths() []string {
	return FindConfigsDownwards("uv.toml")
}

// GetLocalBunfigPaths returns all local '.bunfig.toml' files found.
func GetLocalBunfigPaths() []string {
	return FindConfigsDownwards(".bunfig.toml")
}
