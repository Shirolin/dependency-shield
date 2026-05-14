package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestAuditCommand(t *testing.T) {
	// Disable color for testing to avoid escape sequences in output
	color.NoColor = true

	buf := new(bytes.Buffer)
	SetOut(buf)

	// Set args for rootCmd to execute audit
	rootCmd.SetArgs([]string{"audit"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute audit command: %v", err)
	}

	output := buf.String()
	expected := "🛡️ DependencyShield Audit Report"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", expected, output)
	}
}
