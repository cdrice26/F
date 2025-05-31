package helper

import (
	"os"
	"sync"
)

type workerFunc func(os.FileInfo, string) error

func RunConcurrent(task workerFunc, workerCount int, matches []string) error {
	var wg sync.WaitGroup
	jobs := make(chan string)
	errs := make(chan error, len(matches))

	worker := func() {
		defer wg.Done()
		for match := range jobs {
			info, err := os.Stat(match)
			if err != nil {
				errs <- err
				continue
			}

			err = task(info, match)

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

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
