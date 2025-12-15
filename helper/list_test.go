package helper

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestGetDirectoryFromArgs verifies behavior when args are too short (returns cwd)
// and when an argument is provided (returns that argument).
func TestGetDirectoryFromArgs(t *testing.T) {
	t.Parallel()

	// Case: args too short -> should return current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	dir, err := GetDirectoryFromArgs([]string{}, 1)
	if err != nil {
		t.Fatalf("GetDirectoryFromArgs returned error: %v", err)
	}
	if dir != cwd {
		t.Fatalf("expected cwd %q, got %q", cwd, dir)
	}

	// Case: arg provided -> return that argument
	want := "/tmp/somewhere"
	// On Windows, use drive-friendly path
	if runtime.GOOS == "windows" {
		want = `C:\somewhere`
	}
	dir2, err := GetDirectoryFromArgs([]string{"prog", want}, 2)
	if err != nil {
		t.Fatalf("GetDirectoryFromArgs returned error: %v", err)
	}
	if dir2 != want {
		t.Fatalf("expected %q, got %q", want, dir2)
	}
}

// TestFormatSize checks formatting at bytes/KB/MB/GB boundaries.
func TestFormatSize(t *testing.T) {
	t.Parallel()

	if got := FormatSize(500); !strings.Contains(got, "500") || !strings.Contains(got, "Bytes") {
		t.Fatalf("unexpected FormatSize for 500: %s", got)
	}
	if got := FormatSize(1024); got != "1.00 KB" {
		t.Fatalf("expected 1.00 KB for 1024, got %s", got)
	}
	if got := FormatSize(1024 * 1024); got != "1.00 MB" {
		t.Fatalf("expected 1.00 MB for 1048576, got %s", got)
	}
	if got := FormatSize(1024 * 1024 * 1024); got != "1.00 GB" {
		t.Fatalf("expected 1.00 GB for 1073741824, got %s", got)
	}
}

// TestGetDirSize creates a nested directory with files and ensures sizes are summed.
func TestGetDirSize(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	// create files with known sizes
	f1 := filepath.Join(td, "a.txt")
	f2dir := filepath.Join(td, "sub")
	if err := os.MkdirAll(f2dir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	f2 := filepath.Join(f2dir, "b.txt")

	content1 := []byte("hello")    // 5 bytes
	content2 := []byte("world!!!") // 8 bytes
	if err := os.WriteFile(f1, content1, 0o644); err != nil {
		t.Fatalf("write f1: %v", err)
	}
	if err := os.WriteFile(f2, content2, 0o644); err != nil {
		t.Fatalf("write f2: %v", err)
	}

	size, err := GetDirSize(td)
	if err != nil {
		t.Fatalf("GetDirSize returned error: %v", err)
	}
	if size < int64(len(content1)+len(content2)) {
		t.Fatalf("GetDirSize returned too small size: got %d want at least %d", size, len(content1)+len(content2))
	}
}

// TestIsParentHidden ensures that hidden parent directories are detected.
func TestIsParentHidden(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	// create a hidden parent inside td
	hidden := filepath.Join(td, ".hidden")
	nested := filepath.Join(hidden, "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir hidden nested failed: %v", err)
	}
	// currentPath is nested; basePath is td -> isParentHidden should be true
	if !isParentHidden(nested, td) {
		t.Fatalf("expected isParentHidden to be true for %q under %q", nested, td)
	}
	// when no hidden parents, should be false
	normal := filepath.Join(td, "a", "b")
	if err := os.MkdirAll(normal, 0o755); err != nil {
		t.Fatalf("mkdir normal failed: %v", err)
	}
	if isParentHidden(normal, td) {
		t.Fatalf("expected isParentHidden to be false for %q under %q", normal, td)
	}
}

// TestGetFileListing verifies hidden inclusion/exclusion and entries produced.
func TestGetFileListing(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	visible := filepath.Join(td, "visible.txt")
	hidden := filepath.Join(td, ".hidden.txt")
	if err := os.WriteFile(visible, []byte("v"), 0o644); err != nil {
		t.Fatalf("write visible: %v", err)
	}
	if err := os.WriteFile(hidden, []byte("h"), 0o644); err != nil {
		t.Fatalf("write hidden: %v", err)
	}

	// exclude hidden
	list, err := GetFileListing(td, false)
	if err != nil {
		t.Fatalf("GetFileListing returned error: %v", err)
	}
	foundVisible := false
	foundHidden := false
	for _, e := range list {
		if e.Path == "visible.txt" {
			foundVisible = true
		}
		if e.Path == ".hidden.txt" {
			foundHidden = true
		}
	}
	if !foundVisible {
		t.Fatalf("visible file not listed")
	}
	if foundHidden {
		t.Fatalf("hidden file should not be listed when includeHidden=false")
	}

	// include hidden
	list2, err := GetFileListing(td, true)
	if err != nil {
		t.Fatalf("GetFileListing returned error: %v", err)
	}
	foundHidden = false
	for _, e := range list2 {
		if e.Path == ".hidden.txt" {
			foundHidden = true
		}
	}
	if !foundHidden {
		t.Fatalf("hidden file should be listed when includeHidden=true")
	}
}

// TestGetDirectoryTree builds a nested tree including a hidden directory and verifies output.
func TestGetDirectoryTree(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	// create tree:
	// td/
	//   file.txt
	//   dir/
	//     a.txt
	//   .hiddendir/
	//     secret.txt
	if err := os.WriteFile(filepath.Join(td, "file.txt"), []byte("f"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(td, "dir"), 0o755); err != nil {
		t.Fatalf("mkdir dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(td, "dir", "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatalf("write a.txt: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(td, ".hiddendir"), 0o755); err != nil {
		t.Fatalf("mkdir hiddendir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(td, ".hiddendir", "secret.txt"), []byte("s"), 0o644); err != nil {
		t.Fatalf("write secret: %v", err)
	}

	// exclude hidden
	tree, err := GetDirectoryTree(td, false)
	if err != nil {
		t.Fatalf("GetDirectoryTree returned error: %v", err)
	}
	// Ensure entries do not include the hidden dir or its file
	for _, e := range tree {
		if strings.Contains(e.Path, ".hiddendir") || strings.Contains(e.Path, "secret.txt") {
			t.Fatalf("hidden directory should not appear in tree when includeHidden=false: entry %q", e.Path)
		}
	}

	// include hidden
	tree2, err := GetDirectoryTree(td, true)
	if err != nil {
		t.Fatalf("GetDirectoryTree returned error: %v", err)
	}
	foundHidden := false
	for _, e := range tree2 {
		if strings.Contains(e.Path, ".hiddendir") {
			foundHidden = true
			break
		}
	}
	if !foundHidden {
		t.Fatalf("expected hidden directory to appear when includeHidden=true")
	}
}
