package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestRunDelete_Usage_NoArgs verifies that running the delete command with no arguments displays the usage message.
func TestRunDelete_Usage_NoArgs(t *testing.T) {
	cmd := &cobra.Command{}
	// define the flag so GetBool won't error
	cmd.Flags().BoolP("force", "f", false, "Force")

	out := captureOutput(func() {
		runDelete(cmd, []string{})
	})

	if !contains(out, "Usage: delete") {
		t.Fatalf("expected usage message, got: %q", out)
	}
}

// TestRunDelete_NoMatches verifies that running the delete command with a non-existent pattern displays the no-match message.
func TestRunDelete_NoMatches(t *testing.T) {
	td := t.TempDir()
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("force", "f", false, "Force")

	pattern := filepath.Join(td, "no_such_*")

	out := captureOutput(func() {
		runDelete(cmd, []string{pattern})
	})

	if !contains(out, "No files matched the source pattern") {
		t.Fatalf("expected no-match message, got: %q", out)
	}
}

// TestRunDelete_DeleteFileSuccess verifies that running the delete command with a valid file path deletes the file.
func TestRunDelete_DeleteFileSuccess(t *testing.T) {
	td := t.TempDir()
	filePath := filepath.Join(td, "to_delete.txt")
	if err := os.WriteFile(filePath, []byte("delete me"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("force", "f", true, "Force")

	out := captureOutput(func() {
		runDelete(cmd, []string{filePath})
	})

	if !contains(out, "Deleted") {
		t.Fatalf("expected deleted message, got: %q", out)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed, stat error: %v", err)
	}
}

// TestRunDelete_WildcardDeletesMultipleFiles verifies that running the delete command with a wildcard pattern deletes multiple files.
func TestRunDelete_WildcardDeletesMultipleFiles(t *testing.T) {
	td := t.TempDir()
	pattern := filepath.Join(td, "batch_*.txt")
	names := []string{"batch_a.txt", "batch_b.txt", "batch_c.txt"}
	for _, n := range names {
		path := filepath.Join(td, n)
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatalf("failed to create file %s: %v", path, err)
		}
	}

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("force", "f", true, "Force")

	out := captureOutput(func() {
		runDelete(cmd, []string{pattern})
	})

	// Expect at least one Deleted message
	if !contains(out, "Deleted") {
		t.Fatalf("expected at least one deleted message, got: %q", out)
	}

	// Ensure all files removed
	for _, n := range names {
		path := filepath.Join(td, n)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed, stat error: %v", path, err)
		}
	}
}

// TestRunDelete_ForceFlag_AllowsDeletionOfDirectory verifies that running the delete command with the force flag allows deletion of a directory.
func TestRunDelete_ForceFlag_AllowsDeletionOfDirectory(t *testing.T) {
	td := t.TempDir()
	dir := filepath.Join(td, "mydir")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	// create a file inside the dir so a non-force delete might fail depending on helper implementation
	nested := filepath.Join(dir, "inner.txt")
	if err := os.WriteFile(nested, []byte("inner"), 0o644); err != nil {
		t.Fatalf("failed to create nested file: %v", err)
	}

	cmd := &cobra.Command{}
	// set the force flag true so helper.Delete is called with force=true
	cmd.Flags().BoolP("force", "f", true, "Force")

	out := captureOutput(func() {
		// pass directory path directly
		runDelete(cmd, []string{dir})
	})

	if !contains(out, "Deleted") {
		t.Fatalf("expected deleted message for directory when force=true, got: %q", out)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("expected directory to be removed, stat error: %v", err)
	}
}

// TestRunDelete_InteractiveConfirmation verifies that running the delete command with interactive confirmation deletes the file.
func TestRunDelete_InteractiveConfirmation(t *testing.T) {
	td := t.TempDir()
	filePath := filepath.Join(td, "interactive_delete.txt")
	if err := os.WriteFile(filePath, []byte("delete me"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := &cobra.Command{}
	// leave force as false so the command will prompt
	cmd.Flags().BoolP("force", "f", false, "Force")

	// Replace stdin with a pipe that writes 'y\n' to simulate user confirmation
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe for stdin: %v", err)
	}
	// write confirmation in a goroutine and close writer
	go func() {
		_, _ = w.Write([]byte("y\n"))
		_ = w.Close()
	}()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	out := captureOutput(func() {
		runDelete(cmd, []string{filePath})
	})

	if !contains(out, "Deleted") {
		t.Fatalf("expected deleted message after interactive confirmation, got: %q", out)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed after confirmation, stat error: %v", err)
	}
}
