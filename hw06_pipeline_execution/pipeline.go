package hw06pipelineexecution

import (
	"errors"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var err error

	if len(stages) == 0 {
		return in
	}

	immediatelyStop := make(Bi)
	defer close(immediatelyStop)

	wg := &sync.WaitGroup{}

	for _, stage := range stages {
		wg.Add(1)
		in, err = worker(in, done, stage, wg)
		if err != nil {
			println("error")
			return in
		}
	}

	wg.Wait()

	return in
}

func worker(in In, done In, stage Stage, wg *sync.WaitGroup) (Out, error) {
	var err error
	out := make(Bi)

	go func() {
		wg.Done()

		stageOut := stage(in)

		defer close(out)

		for i := range stageOut {
			select {
			case out <- i:
				println("out <- i")
			case <-done:
				err = errors.New("done")
				println("<-done")
				return
			}
		}
	}()

	return out, err
}
