# Queue

[![Go Report Card](https://goreportcard.com/badge/github.com/robteix/queue)](https://goreportcard.com/report/github.com/robteix/queue)  [![Documentation](https://godoc.org/github.com/robteix/queue?status.svg)](http://godoc.org/github.com/robteix/queue) [![GitHub issues](https://img.shields.io/github/issues/robteix/queue.svg)](https://github.com/robteix/queue/issues) [![license](https://img.shields.io/github/license/robteix/queue.svg?maxAge=2592000)](https://github.com/robteix/queue/LICENSE)

Queue implements a simple FIFO structure in Go with optional expiration.

Simple use case:

```go
var q queue.Queue

q.Put("foo")
q.Put("bar")
val, _ := q.Next()
// val is "foo"
val, _ = q.Next()
// val is "bar"
var, err := q.Next()
// err as no items available
```

You can also wait for an item to be available --

```go
val, err := q.NextWait(30 * time.Second)
if err != nil {
    log.Fatal("timeout waiting for a new item to be queued")
}
log.Println("val", val)
```

Optionally, you can set the queue to expire items, which is useful to prevent a
queue from growing indefinitely --

```go
// create a queue with a TTL of 1s
q := queue.New(queue.WithTTL(time.Second))
q.Put("foo")
time.Sleep(time.Second)
q.Put("bar")
val, _ := q.Next()
// val == "bar", since "foo" was expired and removed from the queue.
```

## License

Copyright (c) 2018 Roberto Selbach Teixeira  <r@rst.sh>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
