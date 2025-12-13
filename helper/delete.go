package helper

import (
	"bufio"
	"fmt"
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

// Delete deletes a file or directory from src. If force is true, it will delete without confirmation.
func Delete(src string, force bool) error {
	matches, err := filepath.Glob(src)
	if err != nil {
		return err
	}

	run := func(info os.FileInfo, match string) error {
		var input string

		if !force {
			reader := bufio.NewReader(os.Stdin)

			fmt.Printf("Are you sure you want to delete %s? (y/n): ", match)
			input, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return err
			}
		}

		if input == "y\n" || input == "Y\n" || force {
			if info.IsDir() {
				err = DeleteDirectory(match)
			} else {
				err = DeleteFile(match)
			}
		} else {
			return fmt.Errorf("deletion of %s cancelled", match)
		}
		return err
	}

	err = RunConcurrent(run, 4, matches)

	if err != nil {
		return err
	}

	return nil
}
