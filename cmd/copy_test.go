package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestRunCopy_Usage_NoArgs verifies that running the copy command with no arguments displays the usage message.
func TestRunCopy_Usage_NoArgs(t *testing.T) {
	cmd := &cobra.Command{}
	// must define the flag so GetBool won't error
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runCopy(cmd, []string{})
	})

	if !contains(out, "Usage: copy") {
		t.Fatalf("expected usage message, got: %q", out)
	}
}

// TestRunCopy_NoMatches verifies that running the copy command with no matches displays the no-match message.
func TestRunCopy_NoMatches(t *testing.T) {
	td := t.TempDir()
	dst := filepath.Join(td, "dst")
	_ = os.MkdirAll(dst, 0o755)

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	pattern := filepath.Join(td, "no_such_*")

	out := captureOutput(func() {
		runCopy(cmd, []string{pattern, dst})
	})

	if !contains(out, "No files matched the source pattern") {
		t.Fatalf("expected no-match message, got: %q", out)
	}
}

// TestRunCopy_CopyFileSuccess verifies that running the copy command with a single file copies it successfully.
func TestRunCopy_CopyFileSuccess(t *testing.T) {
	td := t.TempDir()
	srcDir := filepath.Join(td, "src")
	dstDir := filepath.Join(td, "dst")
	_ = os.MkdirAll(srcDir, 0o755)
	_ = os.MkdirAll(dstDir, 0o755)

	srcFile := filepath.Join(srcDir, "hello.txt")
	content := []byte("hello world")
	if err := os.WriteFile(srcFile, content, 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runCopy(cmd, []string{srcFile, dstDir})
	})

	if !contains(out, "Copied") {
		t.Fatalf("expected copied message, got: %q", out)
	}

	dstFile := filepath.Join(dstDir, "hello.txt")
	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("expected destination file to exist, but got error: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("destination file content mismatch: got %q want %q", string(got), string(content))
	}
}

// TestRunCopy_CopyDirSuccess verifies that running the copy command with a directory copies it successfully.
func TestRunCopy_CopyDirSuccess(t *testing.T) {
	td := t.TempDir()
	srcDir := filepath.Join(td, "src")
	dstDir := filepath.Join(td, "dst")
	_ = os.MkdirAll(srcDir, 0o755)
	_ = os.MkdirAll(dstDir, 0o755)

	srcFile := filepath.Join(srcDir, "hello.txt")
	content := []byte("hello world")
	if err := os.WriteFile(srcFile, content, 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runCopy(cmd, []string{srcDir, dstDir})
	})

	if !contains(out, "Copied") {
		t.Fatalf("expected copied message, got: %q", out)
	}

	dstFile := filepath.Join(dstDir, "src", "hello.txt")
	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("expected destination file to exist, but got error: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("destination file content mismatch: got %q want %q", string(got), string(content))
	}
}
