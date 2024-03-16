package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return errors.New("n must be greater than 0")
	}

	var count int32
	wg := sync.WaitGroup{}
	ch := make(chan Task)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ch {
				if task() != nil {
					atomic.AddInt32(&count, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&count) >= int32(m) {
			break
		}
		ch <- task
	}

	close(ch)
	wg.Wait()

	if count >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
