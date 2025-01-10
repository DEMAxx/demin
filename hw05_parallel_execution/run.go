package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var t [][]Task

	chunkSize := (len(tasks) + n - 1) / n

	for i := 0; i < len(tasks); i += chunkSize {
		end := i + chunkSize

		if end > len(tasks) {
			end = len(tasks)
		}

		t = append(t, tasks[i:end])
	}

	var e error
	ec := 0

	wg := sync.WaitGroup{}

	for _, v := range t {
		sTasks := v

		for _, task := range sTasks {
			wg.Add(1)

			if ec > m {
				return ErrErrorsLimitExceeded
			}

			go func() {
				err := task()
				if err != nil {
					ec++
				}

				wg.Done()
			}()
		}

		wg.Wait()
	}

	return e
}
