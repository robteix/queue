package queue

import (
	"fmt"
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

func ExampleQueue_NextWait() {
	var q Queue
	go func() {
		// wait 1s and queue something
		time.Sleep(time.Second)
		q.Put("foo")
	}()
	fmt.Printf("%d items in queue; waiting\n", q.Len())
	// wait for up to 2s for a new item
	val, err := q.NextWait(2 * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("val:", val)
	// Output:
	// 0 items in queue; waiting
	// val: foo
}

func ExampleQueue_withExpiration() {
	// create a queue with a 1s expiration
	q := New(WithTTL(time.Second))
	q.Put("foo")
	time.Sleep(time.Second)
	// add "bar" and since 1s has ellapsed, "foo" is expired and will therefore
	// be removed from the queue
	q.Put("bar")
	val, err := q.Next()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("val:", val)
	// Output: val: bar
}
