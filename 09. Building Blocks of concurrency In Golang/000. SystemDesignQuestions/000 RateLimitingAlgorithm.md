## Implementing old school token bucket rate limiter 
<details>

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
</details>

### daemon thread is single one ..why we need to put the lock ?

#### Summary of Why Lock Is Needed:
- Shared Resource: Both the daemon thread and the client goroutines modify the same shared variable (possibleTokens).
- Race Conditions: Without locking, simultaneous access to possibleTokens could lead to race conditions, causing incorrect token counts and inconsistent behavior.
- Atomicity: The lock ensures that token bucket operations (checking and modifying possibleTokens) are atomic, preventing data races and ensuring consistent behavior.

## Implementing scalable token bucket rate limiter using redis lua script

<details>

```go
package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis connection setup
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // Change to your Redis server address
})

// Lua script for atomic token consumption
var consumeTokenScript = redis.NewScript(`
	local tokens = redis.call("GET", KEYS[1])
	if not tokens then
		return -1 -- No bucket found
	end
	tokens = tonumber(tokens)
	if tokens > 0 then
		redis.call("DECR", KEYS[1]) -- Consume a token
		return tokens - 1
	else
		return -2 -- No tokens available
	end
`)

// User-specific token bucket filter
type RedisTokenBucket struct {
	Key         string
	MaxTokens   int
	RefillRate  time.Duration
	RefillCount int
	BurstLimit  int
}

// Initialize token bucket for a user
func (tbf *RedisTokenBucket) InitBucket() {
	// Set initial tokens if not already set
	exists, err := rdb.Exists(ctx, tbf.Key).Result()
	if err != nil {
		log.Println("Redis error:", err)
		return
	}
	if exists == 0 {
		rdb.Set(ctx, tbf.Key, tbf.MaxTokens, 0)
	}
	// Start background refill goroutine
	go tbf.refillTokens()
}

// Background job to refill tokens for each user
func (tbf *RedisTokenBucket) refillTokens() {
	for {
		time.Sleep(tbf.RefillRate)
		currentTokens, err := rdb.Get(ctx, tbf.Key).Int()
		if err != nil {
			log.Println("Error fetching tokens:", err)
			continue
		}

		// Refill tokens up to max + burst limit
		if currentTokens < tbf.MaxTokens+tbf.BurstLimit {
			newTokens := min(tbf.MaxTokens+tbf.BurstLimit, currentTokens+tbf.RefillCount)
			rdb.Set(ctx, tbf.Key, newTokens, 0)
		}
	}
}

// Consume a token for a user
func (tbf *RedisTokenBucket) GetToken(userID string) bool {
	result, err := consumeTokenScript.Run(ctx, rdb, []string{tbf.Key}).Int()
	if err != nil {
		log.Println("Redis error:", err)
		return false
	}

	if result >= 0 {
		fmt.Println("User", userID, "granted token at", time.Now(), "- Remaining:", result)
		return true
	}

	fmt.Println("User", userID, "denied - No tokens available at", time.Now())
	return false
}

// Utility function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var wg sync.WaitGroup

	// Different users with different rate limits
	users := map[string]*RedisTokenBucket{
		"user_123": {Key: "token_bucket:user_123", MaxTokens: 10, RefillRate: time.Second, RefillCount: 1, BurstLimit: 5},
		"user_456": {Key: "token_bucket:user_456", MaxTokens: 5, RefillRate: time.Second, RefillCount: 2, BurstLimit: 3},
	}

	// Initialize token buckets for each user
	for _, bucket := range users {
		bucket.InitBucket()
	}

	// Simulate multiple users requesting tokens
	for i := 1; i <= 10; i++ {
		wg.Add(2) // Two users making requests

		go func(i int) {
			defer wg.Done()
			users["user_123"].GetToken("user_123")
		}(i)

		go func(i int) {
			defer wg.Done()
			users["user_456"].GetToken("user_456")
		}(i)
	}

	wg.Wait()
}
```
</details>

## Why Redis Lua Scripting is Good

### 1. Atomicity
Redis executes Lua scripts atomically — the entire script runs as a single operation, without interruption.

💡 This means no race conditions, even when multiple distributed services or threads access the same Redis key concurrently.

✅ Problem Solved:
If multiple nodes try to increment a counter or refill tokens at the same time, without atomic execution, they might:

Over-issue tokens (exceed the limit),

Step on each other’s updates.

### 2. Consistency Across Distributed Nodes
Since all instances share a centralized Redis store, they operate on a single source of truth — the Redis keys.

✅ Problem Solved:
Without a centralized store, each node would need to maintain its own rate counter — leading to inconsistencies (e.g., each node allows 100 requests, but total becomes 100 * N nodes).

### 3. Performance and Latency
Redis is in-memory and extremely fast. Lua scripts run server-side in Redis, avoiding multiple round-trips between app and Redis.

✅ Problem Solved:
Minimizes network chatter — if you did a GET, compute in app, then SET, you'd need 2+ network calls. Lua does it in 1 atomic call.

## Options for Scaling Redis in a Distributed Environment

### Option A: Standalone Redis
Single Pod + PVC.

Easy to deploy.

Can be a single point of failure (SPOF).

Use when:

Your app tolerates occasional Redis downtime.
You use Redis only for rate limiting or caching, not persistent critical data.

### Option B: Redis Sentinel
Provides failover and high availability.

Needs 3 Sentinel replicas and at least 2 Redis nodes.

Redis client automatically discovers the current master.

Kubernetes deployment:
Use Bitnami Helm chart with Sentinel enabled.

### Option C: Redis Cluster
Shards data across multiple Redis nodes.

Provides horizontal scaling.

Higher complexity, but better for very large-scale systems.

Used for:

Massive global rate-limiting needs.

Storing lots of keys with high throughput.

### Reference
https://godoy-lucas-e.medium.com/golang-concurrency-building-a-simple-rate-limiter-token-bucket-algorithm-62de4f389039
https://dev.to/jacktt/implement-rate-limit-in-golang-l4g
