package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

var mutex sync.Mutex
var mutex2 sync.Mutex

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	//var e error
	count := len(tasks)
	//isEnding := false
	errorCount := 0
	jobsCount := 0
	chanJobs := make(chan Task, count)
	//done := make(chan int)

	wg := sync.WaitGroup{}

	defer func() {
		println("defer")
		//if e == nil {
		//	//wg.Wait()
		//	close(chanJobs)
		//}
	}()

	for _, taskJob := range tasks {
		wg.Add(1)

		go func() {
			defer wg.Done()
			chanJobs <- taskJob
		}()
	}

	wg.Wait()

	for {
		for task := range chanJobs {
			println("chanJob", jobsCount)

			if errorCount >= m {
				return ErrErrorsLimitExceeded
			}

			jobsCount++

			err := task()

			if err != nil {
				errorCount++
			}

			if jobsCount >= count {
				return nil
			}
		}

	}
	//return nil
}
