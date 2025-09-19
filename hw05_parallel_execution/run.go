package hw05parallelexecution

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-taskCh:
					if !ok {
						return
					}
					if err := task(); err != nil && m > 0 {
						count := atomic.AddInt64(&errCount, 1)
						if int(count) >= m {
							cancel()
							return
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go resolveTasks(ctx, tasks, taskCh)
	wg.Wait()

	if m > 0 && int(errCount) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func resolveTasks(ctx context.Context, tasks []Task, taskCh chan Task) {
	func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case taskCh <- task:
			case <-ctx.Done():
				return
			}
		}
	}()
}
