package helper

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CopyFile copies a single file from src to dst. If removeSource is true, the source file is deleted after copying.
func CopyFile(src, dst string, removeSource bool) error {
	_, filename := filepath.Split(src)
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(filepath.Join(dst, filename))
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	if removeSource {
		err = os.Remove(src)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyDirectory copies a directory from src to dst. If removeSource is true, the source directory is deleted after copying.
func CopyDirectory(src, dst string, removeSource bool) error {
	src = strings.TrimSuffix(src, string(os.PathSeparator))
	dst = filepath.Join(dst, filepath.Base(src))

	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, src)
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		err = CopyFile(path, filepath.Dir(targetPath), removeSource)
		if err != nil {
			return err
		}

		return nil
	})
}

// Copy handles copying files, directories, and wildcards.
func Copy(src, dst string, removeSource bool) error {
	matches, err := filepath.Glob(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	jobs := make(chan string)
	errs := make(chan error, len(matches))
	workerCount := 4

	worker := func() {
		defer wg.Done()
		for match := range jobs {
			info, err := os.Stat(match)
			if err != nil {
				errs <- err
				continue
			}

			if info.IsDir() {
				err = CopyDirectory(match, dst, removeSource)
			} else {
				err = CopyFile(match, dst, removeSource)
			}

			if err != nil {
				errs <- err
			}
		}
	}

	for range workerCount {
		wg.Add(1)
		go worker()
	}

	go func() {
		for _, match := range matches {
			jobs <- match
		}
		close(jobs)
	}()

	wg.Wait()
	close(errs)

	// Collect errors
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
