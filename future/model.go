package future

import "time"

type payload[T any] struct {
	val T
	err error
}

type Promise[T any] chan payload[T]

func NewPromise[T any]() Promise[T] {
	return make(chan payload[T])
}

func Async[T any](fa func(Promise[T])) Promise[T] {
	f := NewPromise[T]()
	go fa(f)
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
