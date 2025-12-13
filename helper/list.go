package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

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
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize, err
}
