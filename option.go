package queue

import (
	"time"
)

// Option represents an functional parameter for the Queue
type Option interface {
	apply(*Queue)
}

// helper implementation of Option to create new options
type optionFunc func(q *Queue)

func (f optionFunc) apply(q *Queue) {
	f(q)
}

// WithTTL is an optional parameter that makes the queue expirable using the
// provided ttl
func WithTTL(ttl time.Duration) Option {
	return optionFunc(func(q *Queue) { q.expirable = true; q.ttl = ttl })
}
