```go
package main

import (
	"fmt"
	"sync"
	"time"
)

//Rate bucketing algorithm

/*
The token bucket logic is enforced by the rate at which tokens are added to the bucket,
which controls how quickly the goroutines can proceed.

Although 10 goroutines are spawned at once, only one token is granted per second
(or whatever interval you set), ensuring that requests are rate-limited.

The use of sync.Cond allows goroutines to wait until tokens are available,
enforcing the rate limit even under concurrent requests.
*/

type MultiThreadedTokenBucketFilter struct {
	cond *sync.Cond
	lock sync.Mutex

	maxTokens      int
	possibleTokens int
	oneSecond      time.Duration
}

type TokenBucketFilterFactory struct{}

func (f TokenBucketFilterFactory) MakeTokenBucketFilter(capacity int) *MultiThreadedTokenBucketFilter {
	tbf := NewMultiThreadedTokenBucketFilter(capacity)
	tbf.initDaemonThread()

	return tbf
}

func NewMultiThreadedTokenBucketFilter(capacity int) *MultiThreadedTokenBucketFilter {
	return &MultiThreadedTokenBucketFilter{
		maxTokens:      capacity,
		cond:           sync.NewCond(&sync.Mutex{}),
		possibleTokens: 0,
		oneSecond:      time.Second,
	}
}

func (tbf *MultiThreadedTokenBucketFilter) initDaemonThread() {
	go tbf.daemonThread()
}

func (tbf *MultiThreadedTokenBucketFilter) daemonThread() {
	for {
		tbf.cond.L.Lock()
		if tbf.possibleTokens < tbf.maxTokens {
			tbf.possibleTokens += 1
		}

		tbf.cond.Signal()
		tbf.cond.L.Unlock()

		time.Sleep(tbf.oneSecond)
	}
}

func (tbf *MultiThreadedTokenBucketFilter) GetToken(threadName string) {
	tbf.cond.L.Lock()
	for tbf.possibleTokens == 0 {
		tbf.cond.Wait()
	}

	tbf.possibleTokens--
	tbf.cond.L.Unlock()
	fmt.Println("Granting", threadName, "token at", time.Now())
}

func main() {
	var wg sync.WaitGroup
	bucket := TokenBucketFilterFactory{}.MakeTokenBucketFilter(10)

	// Simulate multiple threads requesting tokens
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			threadName := fmt.Sprintf("Thread_%d", i)
			bucket.GetToken(threadName)
		}(i)
	}

	wg.Wait()
}
```
### daemon thread is single one ..why we need to put the lock ?

#### Summary of Why Lock Is Needed:
- Shared Resource: Both the daemon thread and the client goroutines modify the same shared variable (possibleTokens).
- Race Conditions: Without locking, simultaneous access to possibleTokens could lead to race conditions, causing incorrect token counts and inconsistent behavior.
- Atomicity: The lock ensures that token bucket operations (checking and modifying possibleTokens) are atomic, preventing data races and ensuring consistent behavior.


### Reference
https://godoy-lucas-e.medium.com/golang-concurrency-building-a-simple-rate-limiter-token-bucket-algorithm-62de4f389039
https://dev.to/jacktt/implement-rate-limit-in-golang-l4g
