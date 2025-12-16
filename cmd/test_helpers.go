package cmd

import (
	"bytes"
	"io"
	"os"
)

// captureOutput runs the provided function while capturing anything written to stdout.
// It returns the captured output as a string. Tests in this package use this helper
// to avoid duplicating stdout-capture logic across multiple test files.
func captureOutput(f func()) string {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		// If we can't create a pipe, run the function normally and return an empty string.
		f()
		return ""
	}
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Run the function while stdout is redirected.
	f()

	// Close writer and restore original stdout.
	_ = w.Close()
	os.Stdout = orig

	// Return captured output.
	return <-outC
}

// contains reports whether sub is within s. Tests use this small helper to avoid
// importing bytes repeatedly in each test file.
func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}
