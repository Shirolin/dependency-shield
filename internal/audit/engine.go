package audit

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/shiro/dependency-shield/internal/config"
	"github.com/shiro/dependency-shield/internal/model"
)

// AuditNpm scans for 'min-release-age' in the given path using the provided policy.
func AuditNpm(path string, p config.Policy) model.AuditResult {
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

	target := p.MinAgeDays
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "min-release-age") {
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
		result.Message = "min-release-age not found"
	} else {
		result.Message = "min-release-age is less than " + strconv.Itoa(target)
	}

	return result
}

// AuditPnpm scans for 'minimum-release-age' in the given path using the provided policy.
func AuditPnpm(path string, p config.Policy) model.AuditResult {
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

	targetMins := p.PnpmMinAgeMins()
	target, _ := strconv.Atoi(targetMins)
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
		result.Message = "minimum-release-age is less than " + targetMins
	}

	return result
}

// AuditUv scans for 'exclude-newer' in the given TOML file using the provided policy.
func AuditUv(path string, p config.Policy) model.AuditResult {
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

	target := p.UvExcludeNewer()
	// Try [tool.uv] exclude-newer or top-level exclude-newer
	val, ok := getNestedValue(cfg, "tool", "uv", "exclude-newer")
	if !ok {
		val, ok = cfg["exclude-newer"]
	}

	if ok {
		if s, ok := val.(string); ok {
			result.CurrentVal = s
			if s == target {
				result.Status = model.StatusPassed
				result.Message = "Security policy met"
				return result
			}
		}
	}

	if !ok {
		result.Message = "exclude-newer not found"
	} else {
		result.Message = "exclude-newer is not " + target
	}

	return result
}

// AuditBun scans for 'minimumReleaseAge' under '[install]' in the given TOML file using the provided policy.
func AuditBun(path string, p config.Policy) model.AuditResult {
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
		targetSecs := p.BunMinAgeSecs()
		target, _ := strconv.ParseInt(targetSecs, 10, 64)
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
		result.Message = "minimumReleaseAge is less than " + p.BunMinAgeSecs()
	}

	return result
}

// AuditTool runs the appropriate audit for a tool across multiple configuration paths.
func AuditTool(toolName string, paths []string, p config.Policy) []model.AuditResult {
	var results []model.AuditResult
	for _, path := range paths {
		var res model.AuditResult
		switch strings.ToLower(toolName) {
		case "npm":
			res = AuditNpm(path, p)
		case "pnpm":
			res = AuditPnpm(path, p)
		case "uv":
			res = AuditUv(path, p)
		case "bun":
			res = AuditBun(path, p)
		default:
			res = model.AuditResult{
				ToolName:   toolName,
				ConfigPath: path,
				Status:     model.StatusSkip,
				Message:    "Unknown tool: " + toolName,
			}
		}
		results = append(results, res)
	}
	return results
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
