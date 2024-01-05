package server

import "sync"

type Queue[T any] struct {
	values []T
	mutex *sync.Mutex
}

func (q *Queue[T]) Size () int{
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.values)
}
func (q *Queue[T]) Push(value T) {
	q.mutex.Lock()
	q.values = append(q.values, value)
	q.mutex.Unlock()
}

func (q *Queue[T]) PopAll() []T {
	q.mutex.Lock()
	queueBuffer := make([]T, len(q.values))
	copy(queueBuffer, q.values)
	q.values = q.values[:0]
	q.mutex.Unlock()
	return queueBuffer
}

func NewQueue[T any]() *Queue[T]{
	return &Queue[T]{values: []T{}, mutex: &sync.Mutex{}}
}