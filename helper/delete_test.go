package helper

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDeleteFile ensures DeleteFile removes a single file.
func TestDeleteFile(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	f := filepath.Join(td, "file.txt")
	if err := os.WriteFile(f, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := DeleteFile(f); err != nil {
		t.Fatalf("DeleteFile returned error: %v", err)
	}

	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed, stat err: %v", err)
	}
}

// TestDeleteDirectory ensures DeleteDirectory removes a directory tree.
func TestDeleteDirectory(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	nested := filepath.Join(td, "dir", "sub")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	f := filepath.Join(nested, "g.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	if err := DeleteDirectory(filepath.Join(td, "dir")); err != nil {
		t.Fatalf("DeleteDirectory returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(td, "dir")); !os.IsNotExist(err) {
		t.Fatalf("expected directory to be removed, stat err: %v", err)
	}
}

// TestDelete_ForceDeletesMatches ensures Delete with force=true removes matched files and directories.
func TestDelete_ForceDeletesMatches(t *testing.T) {
	t.Parallel()

	td := t.TempDir()

	// create a file
	f := filepath.Join(td, "to-delete.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	// create a directory with a file
	dir := filepath.Join(td, "delme")
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("y"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	pattern := filepath.Join(td, "*")
	if err := Delete(pattern, true); err != nil {
		t.Fatalf("Delete with force returned error: %v", err)
	}

	// both file and directory should be removed
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat err: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("expected dir removed, stat err: %v", err)
	}
}

// TestDelete_ConfirmYesAndNo exercises interactive confirmation behavior.
// It must not be run in parallel because it modifies os.Stdin.
func TestDelete_ConfirmYesAndNo(t *testing.T) {
	// prepare a temporary file to delete
	td := t.TempDir()
	f := filepath.Join(td, "confirm.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	// Save original Stdin and restore at the end.
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Case 1: user answers 'n' -> deletion should be cancelled and file remain.
	r1, w1, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe creation failed: %v", err)
	}
	// write 'n' response and close writer so reader gets EOF after newline
	if _, err := w1.Write([]byte("n\n")); err != nil {
		t.Fatalf("write to pipe failed: %v", err)
	}
	_ = w1.Close()

	os.Stdin = r1
	if err := Delete(f, false); err == nil {
		t.Fatalf("expected deletion to be cancelled (error), got nil")
	}
	// file should still exist
	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("file was removed but should have remained after 'n' response")
		}
		t.Fatalf("unexpected stat error: %v", err)
	}
	_ = r1.Close()

	// Case 2: user answers 'y' -> deletion should proceed and file removed.
	r2, w2, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe creation failed: %v", err)
	}
	if _, err := w2.Write([]byte("y\n")); err != nil {
		t.Fatalf("write to pipe failed: %v", err)
	}
	_ = w2.Close()

	os.Stdin = r2
	if err := Delete(f, false); err != nil {
		t.Fatalf("expected deletion to proceed after 'y', got error: %v", err)
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed after confirmation, stat err: %v", err)
	}
	_ = r2.Close()
}
