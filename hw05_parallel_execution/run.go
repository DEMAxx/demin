package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var mutex sync.Mutex

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	count := len(tasks)
	isEnding := false
	aj, i, ec := 0, 0, 0
	results := make(chan Task, n)

	for {
		if ec > m {
			return ErrErrorsLimitExceeded
		}

		if i >= count && aj == 0 {
			break
		}

		if i < count {

			if aj < n {
				mutex.Lock()
				task := tasks[i]
				i++
				mutex.Unlock()
				go func() {
					if err := task(); err != nil {
						mutex.Lock()
						ec++
						mutex.Unlock()
					}

					results <- task
				}()

				mutex.Lock()
				aj++
				mutex.Unlock()
			}
		} else {
			isEnding = true
		}

		if aj >= n || isEnding {
			<-results
			mutex.Lock()
			aj--
			mutex.Unlock()
		}
	}

	return nil
}
