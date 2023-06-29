package goroutines

import "sync"

// Limiter restricts the number of simultaneous go routines that can be run.
type Limiter struct {
	wg    sync.WaitGroup
	guard chan struct{}
}

// NewLimiter creates a new limiter with the specified maximum go routines allowed.
func NewLimiter(max int) *Limiter {
	return &Limiter{
		wg:    sync.WaitGroup{},
		guard: make(chan struct{}, max),
	}
}

// Do is a convenience function to run a function under the go routine limitation.
// You may choose to use the Increment and Decrement functions individually instead for more control.
func (l *Limiter) Do(fn func()) {
	l.Increment()
	go func() {
		fn()
		l.Decrement()
	}()
}

// Increment increases the count of current go routines being run.
// Call this before creating a constrained go routine.
func (l *Limiter) Increment() {
	l.wg.Add(1)
	l.guard <- struct{}{}
}

// Decrement decreases the count of current go routines being run.
// Call this at the end of the go routine.
func (l *Limiter) Decrement() {
	l.wg.Done()
	<-l.guard
}

// Wait blocks until all go routines have executed.
// Call this to wait for all your go routines to finish before continuing on the main thread.
func (l *Limiter) Wait() {
	l.wg.Wait()
}
