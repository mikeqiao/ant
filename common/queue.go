package common

import (
	"sync"

	queue "gopkg.in/eapache/queue.v1"
)

type Queue struct {
	lock   sync.Mutex
	cond   *sync.Cond
	buffer *queue.Queue
	closed bool
}

func NewQueue() *Queue {
	q := &Queue{
		buffer: queue.New(),
	}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *Queue) Pop() (v interface{}) {
	c := q.cond
	buffer := q.buffer
	q.lock.Lock()
	if 0 == buffer.Length() && !q.closed {
		c.Wait()
	}
	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
	}
	q.lock.Unlock()
	return
}

func (q *Queue) TryPop() (v interface{}, ok bool) {
	buffer := q.buffer
	q.lock.Lock()
	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		ok = true
	} else if q.closed {
		ok = true
	}
	q.lock.Unlock()
	return
}

func (q *Queue) Push(v interface{}) {
	q.lock.Lock()
	if !q.closed {
		q.buffer.Add(v)
		q.cond.Signal()
	}
	q.lock.Unlock()
}

func (q *Queue) Len() (l int) {
	q.lock.Lock()
	l = q.buffer.Length()
	q.lock.Unlock()
	return
}

func (q *Queue) Close() {
	q.lock.Lock()
	if !q.closed {
		q.closed = true
		q.cond.Signal()
	}
	q.lock.Unlock()
}
