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
	var e error
	count := len(tasks)
	isEnding := false
	stop := false
	i, ec := 0, 0
	chanJobs := make(chan Task, n)

	jobsCount := 0
	wg := sync.WaitGroup{}
	defer func() {
		if e == nil {
			wg.Wait()
		}
	}()

	for _, taskJob := range tasks {
		wg.Add(1)
		i++
		go func() {
			defer wg.Done()

			if stop {
				return
			}
			if err := taskJob(); err != nil {
				ec++
			}

			chanJobs <- taskJob

			if i >= n {
				i = 0
				isEnding = true
			}
		}()
	}

	for {
		if isEnding {
			isEnding = false

			for _ = range chanJobs {
				if ec >= m {
					println("error")
					e = ErrErrorsLimitExceeded
					return ErrErrorsLimitExceeded
				}
				jobsCount++
				if jobsCount >= count {
					stop = true
					return nil
				}
			}

			if stop {
				break
			}
		}
	}

	return nil
}
