package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestRunSearch_NotEnoughArgs verifies runSearch returns an error when no args are provided.
func TestRunSearch_NotEnoughArgs(t *testing.T) {
	err := runSearch([]string{}, func(q, d string) ([]string, error) { return nil, nil })
	if err == nil || !strings.Contains(err.Error(), "not enough arguments") {
		t.Fatalf("expected 'not enough arguments' error, got: %v", err)
	}
}

// TestRunSearch_HelperFuncError ensures runSearch propagates errors returned by the helper function.
func TestRunSearch_HelperFuncError(t *testing.T) {
	args := []string{"query", "/some/dir"}
	wantErr := errors.New("helper failure")
	err := runSearch(args, func(q, d string) ([]string, error) {
		return nil, wantErr
	})
	if err == nil || !errors.Is(err, wantErr) {
		t.Fatalf("expected helper error to be returned, got: %v", err)
	}
}

// TestRunSearch_NoResults ensures runSearch returns an error when helper finds no matches.
func TestRunSearch_NoResults(t *testing.T) {
	args := []string{"query", "/some/dir"}
	err := runSearch(args, func(q, d string) ([]string, error) {
		return []string{}, nil
	})
	if err == nil || !strings.Contains(err.Error(), "no files found for search criteria") {
		t.Fatalf("expected 'no files found for search criteria' error, got: %v", err)
	}
}

// TestRunSearch_PrintsResults ensures runSearch prints the results and returns nil on success.
func TestRunSearch_PrintsResults(t *testing.T) {
	args := []string{"q", "/some/dir"}
	out := captureOutput(func() {
		err := runSearch(args, func(q, d string) ([]string, error) {
			return []string{"/path/one", "/path/two"}, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "/path/one") || !strings.Contains(out, "/path/two") {
		t.Fatalf("expected printed results, got: %q", out)
	}
}

// Integration-style tests using the real helper functions:
// create temp files and ensure runNameSearch and runContentSearch find them.
func TestRunNameAndContentSearch_Integration(t *testing.T) {
	td := t.TempDir()

	// create files
	matchName := filepath.Join(td, "matchfile.txt")
	other := filepath.Join(td, "other.txt")

	if err := os.WriteFile(matchName, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("failed to write match file: %v", err)
	}
	if err := os.WriteFile(other, []byte("no-match"), 0o644); err != nil {
		t.Fatalf("failed to write other file: %v", err)
	}

	// run name search
	outName := captureOutput(func() {
		// runNameSearch expects (cmd *cobra.Command, args []string)
		cmd := &cobra.Command{}
		err := runNameSearch(cmd, []string{"match", td})
		if err != nil {
			t.Fatalf("runNameSearch returned error: %v", err)
		}
	})
	if !strings.Contains(outName, matchName) {
		t.Fatalf("expected name search to print %s, got: %q", matchName, outName)
	}

	// run content search
	outContent := captureOutput(func() {
		cmd := &cobra.Command{}
		err := runContentSearch(cmd, []string{"hello", td})
		if err != nil {
			t.Fatalf("runContentSearch returned error: %v", err)
		}
	})
	if !strings.Contains(outContent, matchName) {
		t.Fatalf("expected content search to print %s, got: %q", matchName, outContent)
	}
}
