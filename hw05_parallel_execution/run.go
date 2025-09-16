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
	taskCh := make(chan Task)
	var errCount int64
	wg := sync.WaitGroup{}

	defer func() {
		close(taskCh)
		wg.Wait()
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskCh {
				if task() != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		taskCh <- task

		if m > 0 && int(atomic.LoadInt64(&errCount)) >= m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
