package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	count := len(tasks)
	isEnding := false
	aj, i, ec := 0, 0, 0
	results := make(chan Task, n)
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for {
		if ec > m {
			return ErrErrorsLimitExceeded
		}

		if i >= count && aj == 0 {
			break
		}

		if i < count {
			if aj < n {
				wg.Add(1)
				task := tasks[i]
				i++
				go func() {
					defer wg.Done()

					if err := task(); err != nil {
						ec++
					}

					results <- task
				}()

				aj++
			}
		} else {
			isEnding = true
		}

		if aj >= n || isEnding {
			<-results
			aj--
		}

	}

	return nil
}
