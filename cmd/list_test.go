package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestRunList_EmptyDir_VerifiesHeader ensures that running list on an empty
// directory prints a header that includes "Name".
func TestRunList_EmptyDir_VerifiesHeader(t *testing.T) {
	td := t.TempDir()

	cmd := &cobra.Command{}
	cmd.Flags().BoolP("no-directory-sizes", "n", false, "Do not display directory sizes")
	cmd.Flags().BoolP("tree", "t", false, "Display directory tree")
	cmd.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")

	out := captureOutput(func() {
		runList(cmd, []string{td})
	})

	if !contains(out, "Name") {
		t.Fatalf("expected output to contain 'Name' header; got: %q", out)
	}
}

// TestRunList_TreeMode ensures that the tree flag changes the header and that
// files in nested directories are shown.
func TestRunList_TreeMode(t *testing.T) {
	td := t.TempDir()
	d1 := filepath.Join(td, "d1")
	if err := os.Mkdir(d1, 0o755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	f1 := filepath.Join(d1, "a.txt")
	if err := os.WriteFile(f1, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	cmd := &cobra.Command{}
	// Set tree flag default to true so GetBool returns true.
	cmd.Flags().BoolP("no-directory-sizes", "n", false, "Do not display directory sizes")
	cmd.Flags().BoolP("tree", "t", true, "Display directory tree")
	cmd.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")

	out := captureOutput(func() {
		runList(cmd, []string{td})
	})

	if !contains(out, "Size") {
		t.Fatalf("expected tree output to contain 'Size' header; got: %q", out)
	}
	if !contains(out, "a.txt") {
		t.Fatalf("expected tree output to list nested file 'a.txt'; got: %q", out)
	}
}

// TestRunList_HiddenFlag ensures hidden files are included only when the hidden
// flag is provided.
func TestRunList_HiddenFlag(t *testing.T) {
	td := t.TempDir()
	hidden := filepath.Join(td, ".secret")
	if err := os.WriteFile(hidden, []byte("hidden"), 0o644); err != nil {
		t.Fatalf("failed to create hidden file: %v", err)
	}

	// Without hidden flag
	cmdNoHidden := &cobra.Command{}
	cmdNoHidden.Flags().BoolP("no-directory-sizes", "n", false, "Do not display directory sizes")
	cmdNoHidden.Flags().BoolP("tree", "t", false, "Display directory tree")
	cmdNoHidden.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")

	outNoHidden := captureOutput(func() {
		runList(cmdNoHidden, []string{td})
	})

	if contains(outNoHidden, ".secret") {
		t.Fatalf("did not expect hidden file when hidden flag is false; output: %q", outNoHidden)
	}

	// With hidden flag enabled
	cmdHidden := &cobra.Command{}
	cmdHidden.Flags().BoolP("no-directory-sizes", "n", false, "Do not display directory sizes")
	cmdHidden.Flags().BoolP("tree", "t", false, "Display directory tree")
	// set the hidden flag default to true so GetBool returns true
	cmdHidden.Flags().BoolP("hidden", "a", true, "Include hidden files and directories")

	outHidden := captureOutput(func() {
		runList(cmdHidden, []string{td})
	})

	if !contains(outHidden, ".secret") {
		t.Fatalf("expected hidden file when hidden flag is true; got: %q", outHidden)
	}
}
