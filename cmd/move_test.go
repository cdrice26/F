package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestRunMove_Usage_NoArgs verifies that running the move command with no arguments displays the usage message.
func TestRunMove_Usage_NoArgs(t *testing.T) {
	cmd := &cobra.Command{}
	// must define the flag so GetBool won't error
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runMove(cmd, []string{})
	})

	if !contains(out, "Usage: move") {
		t.Fatalf("expected usage message, got: %q", out)
	}
}

// TestRunMove_NoMatches verifies that running the move command with no matches displays the no-match message.
func TestRunMove_NoMatches(t *testing.T) {
	td := t.TempDir()
	dst := filepath.Join(td, "dst")
	_ = os.MkdirAll(dst, 0o755)

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	pattern := filepath.Join(td, "no_such_*")

	out := captureOutput(func() {
		runMove(cmd, []string{pattern, dst})
	})

	if !contains(out, "No files matched the source pattern") {
		t.Fatalf("expected no-match message, got: %q", out)
	}
}

// TestRunMove_MoveFile verifies that running the move command moves a file from source to destination
// by leveraging the helper.Copy logic (removeSource=true).
func TestRunMove_MoveFile(t *testing.T) {
	td := t.TempDir()

	srcDir := filepath.Join(td, "src")
	dstDir := filepath.Join(td, "dst")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("failed to create dst dir: %v", err)
	}

	srcFile := filepath.Join(srcDir, "move_me.txt")
	content := []byte("move-content")
	if err := os.WriteFile(srcFile, content, 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runMove(cmd, []string{srcFile, dstDir})
	})

	if !contains(out, "Moved") {
		t.Fatalf("expected output to contain 'Moved', got: %q", out)
	}

	// Source should be removed and destination should have the file
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Fatalf("expected source file to be removed, stat error: %v", err)
	}

	dstFile := filepath.Join(dstDir, "move_me.txt")
	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("expected destination file to exist: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("destination content mismatch: got %q want %q", string(got), string(content))
	}
}
