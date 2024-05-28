## Problem Statement
Five silent philosophers sit at a round table with bowls of spaghetti. Forks are placed between each pair of adjacent philosophers.

Each philosopher must alternately think and eat. However, a philosopher can only eat spaghetti when they have both left and right forks. Each fork can be held by only one philosopher and so a philosopher can use the fork only if it is not being used by another philosopher. After an individual philosopher finishes eating, they need to put down both forks so that the forks become available to others. A philosopher can take the fork on their right or the one on their left as they become available, but cannot start eating before getting both forks.

Eating is not limited by the remaining amounts of spaghetti or stomach space; an infinite supply and an infinite demand are assumed.

Design a discipline of behaviour (a concurrent algorithm) such that no philosopher will starve; i.e., each can forever continue to alternate between eating and thinking, assuming that no philosopher can know when others may want to eat or think.


```golang
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type DiningPhilosopher struct {
	forks     []*sync.Mutex
	maxDiners chan struct{}
	exit      chan struct{}
}

func NewDiningPhilosopherObj() *DiningPhilosopher {
	dp := &DiningPhilosopher{
		forks:     make([]*sync.Mutex, 5),
		maxDiners: make(chan struct{}, 4),
		exit:      make(chan struct{}),
	}

	for i := range dp.forks {
		dp.forks[i] = &sync.Mutex{}
	}

	return dp

}

func (dp *DiningPhilosopher) lifecycleOfPhilosopher(id int) {
	for {
		select {
		case <-dp.exit:
			return
		default:
			dp.contemplate()
			dp.eat(id)
		}
	}
}

func (dp *DiningPhilosopher) contemplate() {
	sleepFor := time.Duration(rand.Intn(401)+800) * time.Millisecond
	time.Sleep(sleepFor)
}

func (dp *DiningPhilosopher) eat(id int) {
	dp.maxDiners <- struct{}{}
	firstFork := dp.forks[id]
	secondFork := dp.forks[(id+4)%5]

	firstFork.Lock()
	secondFork.Lock()

	fmt.Printf("Philosopher %d is eating\n", id)

	firstFork.Unlock()
	secondFork.Unlock()

	<-dp.maxDiners
}

func main() {
	dp := NewDiningPhilosopherObj()
	var wg sync.WaitGroup

	/*
		The wg.Add(5) call is used to increment the wait group counter by 5.
		This informs the wait group that there are 5 goroutines that need to complete
		their execution before the main goroutine can proceed past the wg.Wait() call.
		Each time a new goroutine is started, it should call wg.Done() when it completes its work
		to decrement the wait group counter. Once the counter reaches zero,
		meaning all 5 goroutines have finished, the main goroutine can continue execution past the wg.Wait() call.
	*/
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			dp.lifecycleOfPhilosopher(id)
		}(i)
	}

	time.Sleep(6 * time.Second)
	close(dp.exit)
	wg.Wait()

}
```
## Explanation of the above is here
https://github.com/MeSabya/PythonConcurrency-/blob/master/DinningPhilosopher/DiningPhilosopher.md 
