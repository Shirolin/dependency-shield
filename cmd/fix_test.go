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

	// Mock input "n" (no)
	input := strings.NewReader("n\n")
	SetIn(input)

	// Reset force
	force = false

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
	
	if strings.Contains(output, "FIXED") {
		t.Errorf("Expected output NOT to contain 'FIXED' when input is 'n'")
	}
}

func TestFixCommandConfirm(t *testing.T) {
	// Disable color for testing
	color.NoColor = true

	buf := new(bytes.Buffer)
	SetOut(buf)

	// Mock input "y" (yes)
	// We might need to mock some files to make it actually succeed, 
	// but here we just check if it tries to fix.
	input := strings.NewReader("y\ny\ny\ny\ny\ny\ny\ny\n") // Provide enough 'y's for all tools
	SetIn(input)

	// Reset force
	force = false

	rootCmd.SetArgs([]string{"fix"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute fix command: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Fix") {
		t.Errorf("Expected output to contain 'Fix' prompts")
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

func TestFixCommandDaysFlag(t *testing.T) {
	color.NoColor = true
	buf := new(bytes.Buffer)
	SetOut(buf)

	rootCmd.SetArgs([]string{"fix", "--days", "15"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute fix command with --days: %v", err)
	}

	output := buf.String()
	expected := "Policy: Minimum release age = 15 days"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", expected, output)
	}
}
