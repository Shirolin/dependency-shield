package model

type ToolStatus string

const (
	StatusPassed ToolStatus = "PASSED"
	StatusFailed ToolStatus = "FAILED"
	StatusWarn   ToolStatus = "WARNING"
	StatusSkip   ToolStatus = "SKIPPED"
)

type AuditResult struct {
	ToolName    string
	ConfigPath  string
	CurrentVal  string
	Status      ToolStatus
	Message     string
}
