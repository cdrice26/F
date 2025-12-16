package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestRunRename_Usage_NoArgs verifies that running the rename command with no arguments displays the usage message.
func TestRunRename_Usage_NoArgs(t *testing.T) {
	cmd := &cobra.Command{}
	// define the flag so GetBool won't error
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runRename(cmd, []string{})
	})

	if !contains(out, "Usage: rename") {
		t.Fatalf("expected usage message, got: %q", out)
	}
}

// TestRunRename_Success verifies that running the rename command successfully renames the file.
func TestRunRename_Success(t *testing.T) {
	td := t.TempDir()
	src := filepath.Join(td, "old.txt")
	if err := os.WriteFile(src, []byte("original"), 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	cmd := &cobra.Command{}
	// overwrite flag default false is fine since destination does not exist
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runRename(cmd, []string{src, "new.txt"})
	})

	if !contains(out, "Renamed") {
		t.Fatalf("expected Renamed message, got: %q", out)
	}

	newPath := filepath.Join(td, "new.txt")
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected renamed file to exist: %v", err)
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("expected original file to be removed, stat error: %v", err)
	}
}

// TestRunRename_DestinationExists_NoOverwrite verifies that running the rename command with a destination that already exists and overwrite=false fails.
func TestRunRename_DestinationExists_NoOverwrite(t *testing.T) {
	td := t.TempDir()
	src := filepath.Join(td, "a.txt")
	dst := filepath.Join(td, "b.txt")

	if err := os.WriteFile(src, []byte("A"), 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}
	if err := os.WriteFile(dst, []byte("B"), 0o644); err != nil {
		t.Fatalf("failed to write destination file: %v", err)
	}

	cmd := &cobra.Command{}
	// default overwrite=false
	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite")

	out := captureOutput(func() {
		runRename(cmd, []string{src, "b.txt"})
	})

	if !contains(out, "Error checking destination file") {
		t.Fatalf("expected error about existing destination, got: %q", out)
	}

	// Ensure files are unchanged
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("expected source still to exist: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("expected destination still to exist: %v", err)
	}
}

// TestRunRename_DestinationExists_OverwriteTrue verifies that running the rename command with a destination that already exists and overwrite=true succeeds.
func TestRunRename_DestinationExists_OverwriteTrue(t *testing.T) {
	td := t.TempDir()
	src := filepath.Join(td, "x.txt")
	dst := filepath.Join(td, "y.txt")

	if err := os.WriteFile(src, []byte("from-src"), 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}
	if err := os.WriteFile(dst, []byte("from-dst"), 0o644); err != nil {
		t.Fatalf("failed to write destination file: %v", err)
	}

	cmd := &cobra.Command{}
	// set overwrite true so rename proceeds even if destination exists
	cmd.Flags().BoolP("overwrite", "o", true, "Overwrite")

	out := captureOutput(func() {
		runRename(cmd, []string{src, "y.txt"})
	})

	// On Unix-like systems os.Rename will replace the destination; ensure success message.
	if !contains(out, "Renamed") {
		t.Fatalf("expected Renamed message when overwrite=true, got: %q", out)
	}

	// src should be gone, dst (with new name) should exist and contain source content.
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("expected source to be removed after rename, stat err: %v", err)
	}
	newPath := filepath.Join(td, "y.txt")
	b, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("expected destination to exist after overwrite rename: %v", err)
	}
	if string(b) != "from-src" {
		t.Fatalf("destination content mismatch: got %q want %q", string(b), "from-src")
	}
}
