package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	if in == nil {
		out := make(chan interface{})
		close(out)
		return out
	}

	job := jobWithDone(in, done)

	for _, stage := range stages {
		curStage := stage(job)
		job = jobWithDone(curStage, done)
	}

	return job
}

func jobWithDone(in In, done In) Out {
	out := make(chan interface{})

	go func() {
		defer func() {
			close(out)
			for i := range in {
				_ = i
			}
		}()

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return out
}
