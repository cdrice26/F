package helper

import (
	"errors"
	"os"
	"sync/atomic"
	"testing"
)

// TestRunConcurrentSuccess verifies that all provided matches are processed
// successfully by RunConcurrent and that no error is returned.
func TestRunConcurrentSuccess(t *testing.T) {
	t.Parallel()

	// create several temp files
	tdir := t.TempDir()
	var matches []string
	arr := [5]int{1, 2, 3, 4, 5}
	for range arr {
		f, err := os.CreateTemp(tdir, "file-*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		_ = f.Close()
		matches = append(matches, f.Name())
	}

	var processed int32
	task := func(info os.FileInfo, path string) error {
		if info == nil {
			return errors.New("nil info")
		}
		atomic.AddInt32(&processed, 1)
		return nil
	}

	if err := RunConcurrent(task, 3, matches); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if got := int(atomic.LoadInt32(&processed)); got != len(matches) {
		t.Fatalf("expected %d processed, got %d", len(matches), got)
	}
}

// TestRunConcurrentTaskError ensures that if a task returns an error for one
// of the matches, RunConcurrent propagates that error.
func TestRunConcurrentTaskError(t *testing.T) {
	t.Parallel()

	tdir := t.TempDir()
	f1, err := os.CreateTemp(tdir, "a-*.txt")
	if err != nil {
		t.Fatalf("create temp file1: %v", err)
	}
	_ = f1.Close()

	f2, err := os.CreateTemp(tdir, "b-*.txt")
	if err != nil {
		t.Fatalf("create temp file2: %v", err)
	}
	_ = f2.Close()

	matches := []string{f1.Name(), f2.Name()}

	expected := errors.New("task failure")

	task := func(info os.FileInfo, path string) error {
		// return error for the second file
		if path == f2.Name() {
			return expected
		}
		return nil
	}

	err = RunConcurrent(task, 2, matches)
	if err == nil {
		t.Fatalf("expected an error but got nil")
	}
	// The returned error should be the same error instance we returned from the task.
	if !errors.Is(err, expected) && err.Error() != expected.Error() {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}

// TestRunConcurrentStatError verifies that an os.Stat error (e.g. file does not exist)
// is propagated by RunConcurrent.
func TestRunConcurrentStatError(t *testing.T) {
	t.Parallel()

	// Non-existent path
	matches := []string{"/path/that/does/not/exist-or-should-not-exist-123456"}

	task := func(info os.FileInfo, path string) error {
		// should not be called in this case
		t.Fatalf("task should not have been called for path: %s", path)
		return nil
	}

	err := RunConcurrent(task, 1, matches)
	if err == nil {
		t.Fatalf("expected an error from os.Stat but got nil")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got: %v", err)
	}
}
