package cmd

import "testing"

// TestRootCmd_Subcommands ensures rootCmd is configured with the expected subcommands.
func TestRootCmd_Subcommands(t *testing.T) {
	// Validate basic root command properties
	if rootCmd == nil {
		t.Fatalf("rootCmd is nil")
	}
	if rootCmd.Use != "f" {
		t.Fatalf("unexpected rootCmd.Use: %q", rootCmd.Use)
	}

	expected := map[string]bool{
		"copy":   false,
		"move":   false,
		"rename": false,
		"delete": false,
		"list":   false,
		"search": false,
	}

	for _, c := range rootCmd.Commands() {
		if _, ok := expected[c.Name()]; ok {
			expected[c.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected subcommand %q to be present on rootCmd", name)
		}
	}
}
