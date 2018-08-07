package queue

import (
	"testing"
	"time"
)

func TestQueue_Next(t *testing.T) {
	q := &Queue{
		expirable: true,
		ttl:       2 * time.Second,
	}
	for i := 0; i < 5; i++ {
		q.Put(i)
	}

	for i := 0; i < 5; i++ {
		_, err := q.Next()
		if err != nil {
			t.Errorf("Next(%d) returned error %v", i, err)
		}
	}
	item, err := q.Next()
	if err == nil {
		t.Errorf("Next returned %v, wanted error", item)
	}
}

func TestQueue_Expiration(t *testing.T) {
	q := New(WithTTL(500 * time.Millisecond))

	q.Put(1)
	q.Put(2)

	time.Sleep(500 * time.Millisecond)
	q.Put(3)
	// first 2 items should be expired by now
	if q.Len() != 1 {
		t.Errorf("got len %d, want 0", q.Len())
	}
	if item, err := q.Next(); err != nil {
		t.Errorf("Next returned error: %v", err)
	} else if item != 3 {
		t.Errorf("got %d, want 3", item)
	}
}
