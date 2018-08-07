/*
queue implements a first-in-first-out queue.

The API is purposefully very simple. We queue items with Put() and retrive them
with Next() or NextWait()

The empty value of Queue is a valid queue with no expiration.

	var q Queue
	q.Put("foo")
	q.Put("bar")
	v, err := q.Next()

Expiration

The queue can be optionally configured to expire items after a given period of
time. An expirable queue will expire items on demand when accessed.

Here is a trivial example showing how to create and use a queue with 1s
expiration.

    q := queue.New(queue.WithTTL(time.Second))

    // add a value
    q.Put("foo")
    time.Sleep(time.Second)
    q.Put("bar")

    // not that "foo" will have been expired by the previous call, so "bar" is
    // the first non-expired item and that's what will be returned by Next()
    val, err := q.Next()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("val", val)

*/
package queue

import (
	"errors"
	"sync"
	"time"
)

// Queue is a FIFO queue
type Queue struct {
	expirable bool
	ttl       time.Duration

	mu    sync.Mutex
	items []*item

	waiters []chan *item
}

type item struct {
	expire time.Time
	value  interface{}
}

// New initializes a new queue
func New(opts ...Option) *Queue {
	q := &Queue{}
	for _, o := range opts {
		o.apply(q)
	}
	return q
}

// Len returns the number of items currently queued
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.expire()
	return len(q.items)
}

// Next returns the next available item
func (q *Queue) Next() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.next()
}

// gets the next available element
func (q *Queue) next() (interface{}, error) {
	q.expire()
	if len(q.items) > 0 {
		item := q.items[0]
		q.items = q.items[1:]
		return item.value, nil
	}

	return nil, errors.New("no items left")
}

// NextWait returns the next available item. If no item is queued, it waits for
// a new item to be queued.
func (q *Queue) NextWait(timeout time.Duration) (interface{}, error) {
	q.mu.Lock()

	val, err := q.next()
	if err != nil {
		ch := make(chan *item)
		q.waiters = append(q.waiters, ch)
		q.mu.Unlock()
		select {
		case item := <-ch:
			return item.value, nil
		case <-time.After(timeout):
			// remove the channel from waiters
			q.mu.Lock()
			for i := range q.waiters {
				if q.waiters[i] == ch {
					copy(q.waiters[i:], q.waiters[i+1:])
					q.waiters[len(q.waiters)-1] = nil
					q.waiters = q.waiters[:len(q.waiters)-1]
					break
				}
			}
			q.mu.Unlock()
			close(ch)
			return nil, errors.New("timeout waiting for item")
		}
	}

	q.mu.Unlock()
	return val, nil
}

// Put will put the provided value in the queue
func (q *Queue) Put(value interface{}) {
	item := &item{value: value}
	if q.expirable {
		item.expire = time.Now().Add(q.ttl)
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.waiters) > 0 {
		w := q.waiters[0]
		q.waiters = q.waiters[1:]
		w <- item
		return
	}
	q.items = append(q.items, item)
}

// if our queue is expirable, then we remove expired items
func (q *Queue) expire() {
	if q.expirable {
		i := 0
		for {
			if i == len(q.items) {
				break // we're at the end
			}
			item := q.items[i]
			if item.expire.After(time.Now()) {
				break // we found the first non-expired item
			}
			i++
		}
		q.items = q.items[i:]
	}
}
