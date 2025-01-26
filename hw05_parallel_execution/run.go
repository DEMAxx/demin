package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	var e error

	inputCh := make(chan Task)
	errorCh := make(chan error)

	sgnlCh := make(chan struct{})

	wg := &sync.WaitGroup{}

	j := 0

	go func() {
		defer close(inputCh)

		for i := range tasks {
			select {
			case inputCh <- tasks[i]:
			case <-sgnlCh:
				break
			}
		}
	}()

	go func() {
		for i := 0; i < n; i++ {
			wg.Add(1)

			go worker(wg, inputCh, errorCh, sgnlCh)
		}
		wg.Wait()
		close(errorCh)
	}()

	for range errorCh {
		j++
		if j == m {
			close(sgnlCh)
			e = ErrErrorsLimitExceeded
			wg.Wait()
			break
		}
	}

	return e
}

func worker(wg *sync.WaitGroup, inCh <-chan Task, outCh chan<- error, sgnlCh <-chan struct{}) {
	defer wg.Done()

	for task := range inCh {
		err := task()
		if err != nil {
			select {
			case outCh <- err:
			case <-sgnlCh:
				return
			}
		}
	}
}
