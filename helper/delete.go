package helper

import (
	"os"
	"path/filepath"
)

// DeleteFile deletes a file from src.
func DeleteFile(src string) error {
	err := os.Remove(src)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDirectory deletes a directory from src.
func DeleteDirectory(src string) error {
	err := os.RemoveAll(src)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a file or directory from src.
func Delete(src string) error {
	matches, err := filepath.Glob(src)
	if err != nil {
		return err
	}

	run := func(info os.FileInfo, match string) error {
		if info.IsDir() {
			err = DeleteDirectory(match)
		} else {
			err = DeleteFile(match)
		}
		return err
	}

	err = RunConcurrent(run, 4, matches)

	return nil
}
