package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestFixCommand(t *testing.T) {
	// Disable color for testing
	color.NoColor = true

	buf := new(bytes.Buffer)
	SetOut(buf)

	// Set args for rootCmd to execute fix
	rootCmd.SetArgs([]string{"fix"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute fix command: %v", err)
	}

	output := buf.String()
	expected := "🛠️ DependencyShield Fix Report"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", expected, output)
	}
}

func TestFixCommandForce(t *testing.T) {
	// Disable color for testing
	color.NoColor = true

	buf := new(bytes.Buffer)
	SetOut(buf)

	// Reset force to false before test
	force = false

	// Set args for rootCmd to execute fix --force
	rootCmd.SetArgs([]string{"fix", "--force"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute fix --force command: %v", err)
	}

	if !force {
		t.Errorf("Expected 'force' variable to be true after executing with --force flag")
	}
}
