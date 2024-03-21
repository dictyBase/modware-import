package concurrent

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

// Task represents a generic function that takes an input of type I and returns an output of type O and an error.
type Task[I any, O any] func(input I) (O, error)

type TaskWrapper[I any, O any] struct {
	TaskFunc Task[I, O] // The task function
	Input    I          // The input for the task function
}

func RunTasks[I any, O any](
	taskSlice []TaskWrapper[I, O],
	logger *logrus.Entry,
) error {
	if len(taskSlice) == 0 {
		return nil
	}
	results, err := concurrentRun(taskSlice)
	if err != nil {
		return err
	}
	for _, rec := range results {
		logger.Debug(rec)
	}
	return nil
}

func concurrentRun[I any, O any](taskSlice []TaskWrapper[I, O]) ([]O, error) {
	resultCh, errCh := work(taskSlice)
	results := make([]O, 0) // Add slice to collect results
	errSlice := make([]error, 0)

	for len(taskSlice) > 0 {
		select {
		case err := <-errCh:
			if err != nil {
				errSlice = append(errSlice, err)
			}
			taskSlice = taskSlice[:len(taskSlice)-1] // Decrease the slice length by 1 for each error received
		case result, ok := <-resultCh:
			if !ok {
				// If the resultCh is closed, exit from the for loop
				goto FINISH
			}
			results = append(results, result)        // Collect result
			taskSlice = taskSlice[:len(taskSlice)-1] // Decrease the slice length by 1 for each result received
		}
	}
FINISH:
	return results, errors.Join(errSlice...)
}

func work[I any, O any](taskSlice []TaskWrapper[I, O]) (chan O, chan error) {
	resultCh := make(chan O)
	errCh := make(chan error)
	var wg sync.WaitGroup
	// Run each function in a goroutine
	for _, tsk := range taskSlice {
		wg.Add(1)
		go func(trun TaskWrapper[I, O]) {
			defer wg.Done()
			result, err := trun.TaskFunc(trun.Input)
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- result
		}(tsk)
	}
	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()
	return resultCh, errCh
}
