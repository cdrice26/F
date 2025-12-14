package helper

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SearchByName searches for files by name in a directory.
func SearchByName(query string, dir string) ([]string, error) {
	var results []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.Contains(d.Name(), query) {
			results = append(results, path)
		}
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.New("not found")
	}
	return results, err
}

// SearchByContent searches for files by content in a directory.
func SearchByContent(query string, dir string) ([]string, error) {
	var results []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(content), query) {
			results = append(results, path)
		}
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return nil, errors.New("not found")
	}
	return results, err
}
