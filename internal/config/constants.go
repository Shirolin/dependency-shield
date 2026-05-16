package config

import (
	"fmt"
	"strconv"
)

const (
	DefaultMinAgeDays = 30
	// 环境预设天数
	EnvLocalDays = 7
	EnvCIDays    = 15
	EnvProdDays  = 30
)

// Environment presets
const (
	EnvLocal = "local"
	EnvCI    = "ci"
	EnvProd  = "prod"
)

// Policy represents the security policy configuration.
type Policy struct {
	MinAgeDays int
}

// NewDefaultPolicy returns a policy with default values.
func NewDefaultPolicy() Policy {
	return Policy{MinAgeDays: DefaultMinAgeDays}
}

// NpmMinAge returns the string value for npm's min-release-age.
func (p Policy) NpmMinAge() string {
	return strconv.Itoa(p.MinAgeDays)
}

// PnpmMinAgeMins returns the string value for pnpm's minimum-release-age in minutes.
func (p Policy) PnpmMinAgeMins() string {
	return strconv.Itoa(p.MinAgeDays * 24 * 60)
}

// BunMinAgeSecs returns the string value for bun's minimumReleaseAge in seconds.
func (p Policy) BunMinAgeSecs() string {
	return strconv.FormatInt(int64(p.MinAgeDays)*24*3600, 10)
}

// UvExcludeNewer returns the string value for uv's exclude-newer.
func (p Policy) UvExcludeNewer() string {
	return fmt.Sprintf("%dd", p.MinAgeDays)
}
