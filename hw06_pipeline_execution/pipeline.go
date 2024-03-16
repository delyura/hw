package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		return nil
	}

	for i := range stages {
		ch := make(Bi)
		go func(ch Bi, in In) {
			defer close(ch)
			for {
				select {
				case <-done:
					return
				case value, ok := <-in:
					if !ok {
						return
					}
					ch <- value
				}
			}
		}(ch, in)
		in = stages[i](ch)
	}
	return in
}
