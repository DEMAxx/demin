package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = worker(done, stage(in))
	}

	return in
}

func worker(done In, stageWork Out) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case value, ok := <-stageWork:
				if !ok {
					return
				}
				out <- value
			case <-done:
				return
			}
		}
	}()

	return out
}
