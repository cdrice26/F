package helper

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	Path     string
	DirEntry fs.DirEntry
}

// GetDirectoryFromArgs returns the directory path from the command line arguments.
func GetDirectoryFromArgs(args []string) (string, error) {
	if len(args) == 0 {
		return os.Getwd()
	} else {
		return args[0], nil
	}
}

// FormatSize formats the size in bytes to a human-readable format.
func FormatSize(size int64) string {
	if size >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(1024*1024*1024))
	} else if size >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(1024*1024))
	} else if size >= 1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(1024))
	}
	return fmt.Sprintf("%d Bytes", size)
}

// Calculate total size of a directory (recursively)
func GetDirSize(path string) (int64, error) {
	var totalSize int64
	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize, err
}

// GetFileListing returns a list of files in a directory.
func GetFileListing(path string) ([]Entry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, file := range files {
		fullPath := file.Name()
		entries = append(entries, Entry{Path: fullPath, DirEntry: file})
	}
	return entries, nil
}

// GetDirectoryTree returns a tree structure of a directory.
func GetDirectoryTree(path string) ([]Entry, error) {
	var entries []Entry
	err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == currentPath {
			entries = append(entries, Entry{Path: path, DirEntry: d})
		} else {
			relPath := strings.Replace(currentPath, path, "", 1)
			nestCount := strings.Count(relPath, string(filepath.Separator))
			pathStr := strings.Repeat("│   ", nestCount)
			fullPath := filepath.Join(pathStr, d.Name())
			entries = append(entries, Entry{Path: fullPath, DirEntry: d})
		}
		return nil
	})

	return entries, err
}
