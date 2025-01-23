package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	inputCh := make(chan Task)
	errorCh := make(chan error)

	wg := &sync.WaitGroup{}

	go func() {
		defer close(inputCh)

		for i := range tasks {
			inputCh <- tasks[i]
		}
	}()

	go func() {
		for i := 0; i < n; i++ {
			wg.Add(1)

			go worker(wg, inputCh, errorCh)
		}
		wg.Wait()
		close(errorCh)
	}()

	j := 0

	for _ = range errorCh {
		j++
		if j == m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func worker(wg *sync.WaitGroup, inCh <-chan Task, outCh chan<- error) {
	defer wg.Done()

	for task := range inCh {
		err := task()

		if err != nil {
			outCh <- err
		}
	}
}
