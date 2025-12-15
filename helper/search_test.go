package helper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSearchByName and TestSearchByContent create files and ensure searches find expected files.
func TestSearchByNameAndContent(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	// files:
	// td/alpha.txt (contains "needle")
	// td/beta.txt  (contains "nothing")
	// td/sub/gamma.txt (contains "needle")
	if err := os.WriteFile(filepath.Join(td, "alpha.txt"), []byte("needle here"), 0o644); err != nil {
		t.Fatalf("write alpha: %v", err)
	}
	if err := os.WriteFile(filepath.Join(td, "beta.txt"), []byte("nothing"), 0o644); err != nil {
		t.Fatalf("write beta: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(td, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.WriteFile(filepath.Join(td, "sub", "gamma.txt"), []byte("needle again"), 0o644); err != nil {
		t.Fatalf("write gamma: %v", err)
	}

	// Search by name: query "alpha" should find alpha.txt
	nameResults, err := SearchByName("alpha", td)
	if err != nil {
		t.Fatalf("SearchByName returned error: %v", err)
	}
	if len(nameResults) != 1 {
		t.Fatalf("SearchByName expected 1 result, got %d: %v", len(nameResults), nameResults)
	}
	if !strings.HasSuffix(nameResults[0], "alpha.txt") {
		t.Fatalf("SearchByName returned unexpected path: %s", nameResults[0])
	}

	// Search by content: query "needle" should find alpha and gamma
	contentResults, err := SearchByContent("needle", td)
	if err != nil {
		t.Fatalf("SearchByContent returned error: %v", err)
	}
	if len(contentResults) != 2 {
		t.Fatalf("SearchByContent expected 2 results, got %d: %v", len(contentResults), contentResults)
	}

	// Search for a non-existing name/content -> expect empty results and no error
	nres, err := SearchByName("does-not-exist", td)
	if err != nil {
		t.Fatalf("SearchByName unexpected error for no matches: %v", err)
	}
	if len(nres) != 0 {
		t.Fatalf("expected zero results for SearchByName no-match, got %d", len(nres))
	}
	cres, err := SearchByContent("nope-nope", td)
	if err != nil {
		t.Fatalf("SearchByContent unexpected error for no matches: %v", err)
	}
	if len(cres) != 0 {
		t.Fatalf("expected zero results for SearchByContent no-match, got %d", len(cres))
	}
}
