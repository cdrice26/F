package helper

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile_OverwriteAndRemoveSource(t *testing.T) {
	tdir := t.TempDir()
	srcDir := filepath.Join(tdir, "src")
	dstDir := filepath.Join(tdir, "dst")

	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("mkdir dst: %v", err)
	}

	// create initial source file
	srcPath := filepath.Join(srcDir, "hello.txt")
	original := []byte("original content")
	if err := os.WriteFile(srcPath, original, 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	// first copy should succeed
	if err := CopyFile(srcPath, dstDir, false, false); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// verify destination content
	dstPath := filepath.Join(dstDir, "hello.txt")
	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("read dest: %v", err)
	}
	if string(got) != string(original) {
		t.Fatalf("dest content mismatch: got %q want %q", string(got), string(original))
	}

	// create another source with same filename but different content
	srcPath2 := filepath.Join(srcDir, "hello_overwrite.txt")
	// We'll use same basename as existing destination to exercise overwrite logic:
	// create a file with the same base name as dstPath but in a different source file.
	// To do that, create a new temp file then rename it to have same basename.
	if err := os.WriteFile(srcPath2, []byte("new content"), 0o644); err != nil {
		t.Fatalf("write source2: %v", err)
	}
	// rename to same basename as dst (simulate same basename)
	srcPath2Renamed := filepath.Join(srcDir, "hello.txt")
	if err := os.Rename(srcPath2, srcPath2Renamed); err != nil {
		t.Fatalf("rename source2: %v", err)
	}

	// attempt to copy without overwrite should return an error
	if err := CopyFile(srcPath2Renamed, dstDir, false, false); err == nil {
		t.Fatalf("expected error copying file when destination exists and overwrite=false")
	}

	// now copy with overwrite=true should succeed and replace content
	newContent := []byte("overwritten content")
	if err := os.WriteFile(srcPath2Renamed, newContent, 0o644); err != nil {
		t.Fatalf("write new content: %v", err)
	}
	if err := CopyFile(srcPath2Renamed, dstDir, false, true); err != nil {
		t.Fatalf("CopyFile with overwrite failed: %v", err)
	}
	got2, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("read dest after overwrite: %v", err)
	}
	if string(got2) != string(newContent) {
		t.Fatalf("dest content mismatch after overwrite: got %q want %q", string(got2), string(newContent))
	}

	// test removeSource true - create a fresh source file and copy+remove it
	srcRemove := filepath.Join(srcDir, "toremove.txt")
	removeContent := []byte("remove me")
	if err := os.WriteFile(srcRemove, removeContent, 0o644); err != nil {
		t.Fatalf("write srcRemove: %v", err)
	}
	if err := CopyFile(srcRemove, dstDir, true, false); err != nil {
		t.Fatalf("CopyFile with removeSource failed: %v", err)
	}
	// source should be removed
	if _, err := os.Stat(srcRemove); !os.IsNotExist(err) {
		t.Fatalf("expected source to be removed, stat err: %v", err)
	}
	// destination should exist
	if _, err := os.Stat(filepath.Join(dstDir, "toremove.txt")); err != nil {
		t.Fatalf("expected destination to exist after copy, stat err: %v", err)
	}
}

func TestCopyDirectory_BasicAndRemoveSource(t *testing.T) {
	tdir := t.TempDir()
	// create a sample directory tree:
	// srcroot/
	//   a.txt
	//   nested/
	//     b.txt
	srcRoot := filepath.Join(tdir, "srcroot")
	if err := os.MkdirAll(filepath.Join(srcRoot, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir tree: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcRoot, "a.txt"), []byte("A"), 0o644); err != nil {
		t.Fatalf("write a.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcRoot, "nested", "b.txt"), []byte("B"), 0o644); err != nil {
		t.Fatalf("write b.txt: %v", err)
	}

	dstParent := filepath.Join(tdir, "dstparent")
	if err := os.MkdirAll(dstParent, 0o755); err != nil {
		t.Fatalf("mkdir dstparent: %v", err)
	}

	// copy directory (without removing source)
	if err := CopyDirectory(srcRoot, dstParent, false, false); err != nil {
		t.Fatalf("CopyDirectory failed: %v", err)
	}

	dstRoot := filepath.Join(dstParent, filepath.Base(srcRoot))
	// check files exist in destination
	aDst := filepath.Join(dstRoot, "a.txt")
	bDst := filepath.Join(dstRoot, "nested", "b.txt")
	if data, err := os.ReadFile(aDst); err != nil {
		t.Fatalf("read aDst: %v", err)
	} else if string(data) != "A" {
		t.Fatalf("aDst content mismatch: %q", string(data))
	}
	if data, err := os.ReadFile(bDst); err != nil {
		t.Fatalf("read bDst: %v", err)
	} else if string(data) != "B" {
		t.Fatalf("bDst content mismatch: %q", string(data))
	}

	// now test removeSource = true
	srcRoot2 := filepath.Join(tdir, "srcroot2")
	if err := os.MkdirAll(filepath.Join(srcRoot2, "nested2"), 0o755); err != nil {
		t.Fatalf("mkdir tree2: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcRoot2, "c.txt"), []byte("C"), 0o644); err != nil {
		t.Fatalf("write c.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcRoot2, "nested2", "d.txt"), []byte("D"), 0o644); err != nil {
		t.Fatalf("write d.txt: %v", err)
	}

	dstParent2 := filepath.Join(tdir, "dstparent2")
	if err := os.MkdirAll(dstParent2, 0o755); err != nil {
		t.Fatalf("mkdir dstparent2: %v", err)
	}

	if err := CopyDirectory(srcRoot2, dstParent2, true, false); err != nil {
		t.Fatalf("CopyDirectory with removeSource failed: %v", err)
	}
	// after removal, srcRoot2 should either not exist or be empty; expect it to not exist
	if _, err := os.Stat(srcRoot2); !os.IsNotExist(err) {
		t.Fatalf("expected source root2 to be removed, stat err: %v", err)
	}
	// destination should have files
	cDst := filepath.Join(dstParent2, filepath.Base(srcRoot2), "c.txt")
	if data, err := os.ReadFile(cDst); err != nil {
		t.Fatalf("read cDst: %v", err)
	} else if string(data) != "C" {
		t.Fatalf("cDst content mismatch: %q", string(data))
	}
}

func TestCopy_WildcardFilesAndDirectories(t *testing.T) {
	tdir := t.TempDir()
	srcParent := filepath.Join(tdir, "wildsrc")
	if err := os.MkdirAll(filepath.Join(srcParent, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir wildsrc: %v", err)
	}

	// create multiple top-level entries under srcParent
	if err := os.WriteFile(filepath.Join(srcParent, "f1.txt"), []byte("one"), 0o644); err != nil {
		t.Fatalf("write f1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcParent, "f2.txt"), []byte("two"), 0o644); err != nil {
		t.Fatalf("write f2: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcParent, "sub", "sf.txt"), []byte("subfile"), 0o644); err != nil {
		t.Fatalf("write subfile: %v", err)
	}

	dst := filepath.Join(tdir, "wilddst")
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatalf("mkdir dst: %v", err)
	}

	pattern := filepath.Join(srcParent, "*")
	if err := Copy(pattern, dst, false, false); err != nil {
		t.Fatalf("Copy with wildcard failed: %v", err)
	}

	// verify files copied
	if data, err := os.ReadFile(filepath.Join(dst, "f1.txt")); err != nil {
		t.Fatalf("read copied f1: %v", err)
	} else if string(data) != "one" {
		t.Fatalf("copied f1 content mismatch: %q", string(data))
	}
	if data, err := os.ReadFile(filepath.Join(dst, "f2.txt")); err != nil {
		t.Fatalf("read copied f2: %v", err)
	} else if string(data) != "two" {
		t.Fatalf("copied f2 content mismatch: %q", string(data))
	}
	// sub dir should have been copied as a directory entry inside dst (CopyDirectory creates a folder with base name)
	// Depending on copy implementation, sub may appear under dst/sub
	if data, err := os.ReadFile(filepath.Join(dst, "sub", "sf.txt")); err != nil {
		// some implementations might have created dst/sub/sf.txt, if not, try dst/sub/sf.txt under nested folder name
		t.Fatalf("read copied sub/sf.txt: %v", err)
	} else if string(data) != "subfile" {
		t.Fatalf("copied subfile content mismatch: %q", string(data))
	}
}
