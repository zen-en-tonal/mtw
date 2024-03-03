package sync

import (
	"errors"
	"sync"
)

// TryAll runs functions that may fail asynchronously.
// Returns an error immediately if execution of at least one function fails.
func TryAll[T any](arg T, funcs ...func(T) error) error {
	var wg sync.WaitGroup
	wg.Add(len(funcs))

	wc := make(chan struct{})
	ec := make(chan error, len(funcs))
	go func() {
		for _, f := range funcs {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				if f == nil {
					return
				}
				if err := f(arg); err != nil {
					ec <- err
				}
			}(&wg)
		}
		wg.Wait()
		close(wc)
		close(ec)
	}()

	select {
	case err := <-ec:
		return err
	case <-wc:
		return nil
	}
}

// TryAll runs functions that may fail asynchronously.
func TrySome[T any](arg T, funcs ...func(T) error) error {
	var wg sync.WaitGroup
	wg.Add(len(funcs))

	ec := make(chan error, len(funcs))
	for _, f := range funcs {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			if f == nil {
				return
			}
			if err := f(arg); err != nil {
				ec <- err
			}
		}(&wg)
	}
	wg.Wait()
	close(ec)

	var errs []error
	for e := range ec {
		if e != nil {
			errs = append(errs, e)
		}
	}
	return errors.Join(errs...)
}
