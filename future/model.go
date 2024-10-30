package future

import (
	"time"
)

const zeroDuration = 0 * time.Nanosecond

type payload[T any] struct {
	val T
	err error
}

type Promise[T any] chan payload[T]

type Optional[T any] struct {
	Present bool
	Value   T
	Err     error
}

func NewPromise[T any]() Promise[T] {
	return make(chan payload[T])
}

func Async[T any](provider func() (T, error)) Promise[T] {
	f := NewPromise[T]()
	go func() {
		f.Resolve(provider())
	}()
	return f
}

func (f Promise[T]) sendAndClose(p payload[T]) {
	f <- p
	close(f)
}

func (f Promise[T]) Resolve(value T, err error) {
	f.sendAndClose(payload[T]{val: value, err: err})
}

func (f Promise[T]) Value(value T) {
	f.sendAndClose(payload[T]{val: value})
}

func (f Promise[T]) Error(err error) {
	f.sendAndClose(payload[T]{err: err})
}

func (f Promise[T]) Await() (T, error) {
	payload := <-f
	return payload.val, payload.err
}

func (f Promise[T]) AwaitWithTimeout(t time.Duration) (T, error, bool) {
	select {
	case payload := <-f:
		return payload.val, payload.err, true
	case <-time.After(t):
		var zero T
		return zero, nil, false
	}
}

func (f Promise[T]) Get() Optional[T] {
	val, err, present := f.AwaitWithTimeout(zeroDuration)

	return Optional[T]{
		Present: present,
		Value:   val,
		Err:     err,
	}
}

func (f Promise[T]) Then(consumer func(T, error)) {
	go func() {
		consumer(f.Await())
	}()
}
