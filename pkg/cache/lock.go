package cache

import "sync"

// Lock implements cache stampede mitigation via locking.
type Lock struct {
	mutex   sync.Mutex
	results map[string]*LockResult
}

// Attempts to acquire a lock on the given key. If successful, f is invoked and the result is made
// available via the returned LockResult. If unsuccessful, f is not invoked, but the result of the
// successful thread's invocation will be made available via the returned LockResult.
func (l *Lock) Acquire(key string, f func() (interface{}, error)) *LockResult {
	l.mutex.Lock()

	if l.results == nil {
		l.results = map[string]*LockResult{}
	}

	if ret, ok := l.results[key]; ok {
		l.mutex.Unlock()
		return ret
	}

	ret := &LockResult{
		done: make(chan struct{}),
	}
	l.results[key] = ret
	l.mutex.Unlock()

	ret.value, ret.err = f()
	close(ret.done)

	l.mutex.Lock()
	delete(l.results, key)
	l.mutex.Unlock()

	return ret
}

type LockResult struct {
	done  chan struct{}
	value interface{}
	err   error
}

// Returns a channel that gets closed when the result has an available value or error.
func (r *LockResult) Done() <-chan struct{} {
	return r.done
}

// The value computed by the provided function. It is unsafe to invoke this before r.Done() is
// closed.
func (r *LockResult) Value() interface{} {
	return r.value
}

// The error computed by the provided function. It is unsafe to invoke this before r.Done() is
// closed.
func (r *LockResult) Err() error {
	return r.err
}

// Waits for the result to be available and returns the value and error.
func (r *LockResult) Wait() (interface{}, error) {
	<-r.done
	return r.value, r.err
}
