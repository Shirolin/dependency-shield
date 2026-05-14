package audit

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/shiro/dependencyshield/internal/config"
	"github.com/shiro/dependencyshield/internal/model"
)

// AuditNpm scans for 'min-release-age=30' in the given path.
func AuditNpm(path string) model.AuditResult {
	result := model.AuditResult{
		ToolName:   "npm",
		ConfigPath: path,
		Status:     model.StatusFailed,
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			result.Status = model.StatusSkip
			result.Message = "Configuration file not found"
			return result
		}
		result.Message = "Error opening file: " + err.Error()
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "min-release-age") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := strings.TrimSpace(parts[1])
				result.CurrentVal = val
				if val == config.NpmMinAge {
					result.Status = model.StatusPassed
					result.Message = "Security policy met"
					return result
				}
			}
		}
	}

	if result.CurrentVal == "" {
		result.Message = "min-release-age not found"
	} else {
		result.Message = "min-release-age is not " + config.NpmMinAge
	}

	return result
}

// AuditPnpm scans for 'minimum-release-age=43200' (or higher) in the given path.
func AuditPnpm(path string) model.AuditResult {
	result := model.AuditResult{
		ToolName:   "pnpm",
		ConfigPath: path,
		Status:     model.StatusFailed,
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			result.Status = model.StatusSkip
			result.Message = "Configuration file not found"
			return result
		}
		result.Message = "Error opening file: " + err.Error()
		return result
	}
	defer file.Close()

	target, _ := strconv.Atoi(config.PnpmMinAgeMins)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "minimum-release-age") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				valStr := strings.TrimSpace(parts[1])
				result.CurrentVal = valStr
				val, err := strconv.Atoi(valStr)
				if err == nil && val >= target {
					result.Status = model.StatusPassed
					result.Message = "Security policy met"
					return result
				}
			}
		}
	}

	if result.CurrentVal == "" {
		result.Message = "minimum-release-age not found"
	} else {
		result.Message = "minimum-release-age is less than " + config.PnpmMinAgeMins
	}

	return result
}

// AuditUv scans for 'exclude-newer = "30d"' in the given TOML file.
func AuditUv(path string) model.AuditResult {
	result := model.AuditResult{
		ToolName:   "uv",
		ConfigPath: path,
		Status:     model.StatusFailed,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			result.Status = model.StatusSkip
			result.Message = "Configuration file not found"
			return result
		}
		result.Message = "Error reading file: " + err.Error()
		return result
	}

	var cfg map[string]interface{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		result.Message = "Error parsing TOML: " + err.Error()
		return result
	}

	// Try [tool.uv] exclude-newer or top-level exclude-newer
	val, ok := getNestedValue(cfg, "tool", "uv", "exclude-newer")
	if !ok {
		val, ok = cfg["exclude-newer"]
	}

	if ok {
		if s, ok := val.(string); ok {
			result.CurrentVal = s
			if s == config.UvExcludeNewer {
				result.Status = model.StatusPassed
				result.Message = "Security policy met"
				return result
			}
		}
	}

	if !ok {
		result.Message = "exclude-newer not found"
	} else {
		result.Message = "exclude-newer is not " + config.UvExcludeNewer
	}

	return result
}

// AuditBun scans for 'minimumReleaseAge = 2592000' under '[install]' in the given TOML file.
func AuditBun(path string) model.AuditResult {
	result := model.AuditResult{
		ToolName:   "bun",
		ConfigPath: path,
		Status:     model.StatusFailed,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			result.Status = model.StatusSkip
			result.Message = "Configuration file not found"
			return result
		}
		result.Message = "Error reading file: " + err.Error()
		return result
	}

	var cfg map[string]interface{}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		result.Message = "Error parsing TOML: " + err.Error()
		return result
	}

	val, ok := getNestedValue(cfg, "install", "minimumReleaseAge")
	if ok {
		target, _ := strconv.ParseInt(config.BunMinAgeSecs, 10, 64)
		var currentVal int64
		var isValid bool

		switch v := val.(type) {
		case int64:
			currentVal = v
			isValid = true
		case float64:
			currentVal = int64(v)
			isValid = true
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err == nil {
				currentVal = iv
				isValid = true
			}
		}

		if isValid {
			result.CurrentVal = strconv.FormatInt(currentVal, 10)
			if currentVal >= target {
				result.Status = model.StatusPassed
				result.Message = "Security policy met"
				return result
			}
		}
	}

	if !ok {
		result.Message = "minimumReleaseAge not found"
	} else {
		result.Message = "minimumReleaseAge is less than " + config.BunMinAgeSecs
	}

	return result
}

func getNestedValue(m map[string]interface{}, keys ...string) (interface{}, bool) {
	var current interface{} = m
	for _, key := range keys {
		currMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		val, ok := currMap[key]
		if !ok {
			return nil, false
		}
		current = val
	}
	return current, true
}
