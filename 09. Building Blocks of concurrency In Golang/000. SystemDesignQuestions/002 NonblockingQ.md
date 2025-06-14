```python
import threading
import time
import random


class NonBlockingQ:
    def __init__(self, max_size):
        self.q = []
        self.max_size = max_size
        self.lock = threading.Lock()
        self.q_waiting_puts = []  # List of threading.Condition objects
        self.q_waiting_gets = []

    def enqueue(self, val):
        with self.lock:
            if len(self.q) == self.max_size:
                print("Queue is Full")
                cond = threading.Condition()
                self.q_waiting_puts.append((cond, val))
                return cond
            else:
                self.q.append(val)
                if self.q_waiting_gets:
                    cond = self.q_waiting_gets.pop(0)
                    with cond:
                        cond.notify()
                return None

    def dequeue(self):
        with self.lock:
            if self.q:
                item = self.q.pop(0)
                if self.q_waiting_puts:
                    cond, val = self.q_waiting_puts.pop(0)
                    with cond:
                        self.q.append(val)
                        cond.notify()
                return item, None
            else:
                cond = threading.Condition()
                self.q_waiting_gets.append(cond)
                return None, cond


def retry_enqueue(cond, val, q):
    with cond:
        cond.wait()
    print("Retry Enqueue Invoked")
    new_cond = q.enqueue(val)
    if new_cond:
        threading.Thread(target=retry_enqueue, args=(new_cond, val, q)).start()
    else:
        print(f"Item {val} successfully added on a retry")


def retry_dequeue(cond):
    with cond:
        cond.wait()
    print("Retry Dequeue Invoked, consumed item")


def producer(q):
    item = 1
    while True:
        cond = q.enqueue(item)
        if cond:
            threading.Thread(target=retry_enqueue, args=(cond, item, q)).start()
        else:
            print(f"Producer produced item {item}")
        item += 1
        time.sleep(random.randint(1, 3))


def consumer(q):
    while True:
        item, cond = q.dequeue()
        if cond:
            threading.Thread(target=retry_dequeue, args=(cond,)).start()
        else:
            print(f"Consumer consumed item {item}")
        time.sleep(1)


if __name__ == "__main__":
    q = NonBlockingQ(2)

    producer_thread = threading.Thread(target=producer, args=(q,))
    consumer_thread = threading.Thread(target=consumer, args=(q,))

    producer_thread.start()
    consumer_thread.start()

    producer_thread.join()
    consumer_thread.join()
```

```golang
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type NonBlockingQ struct {
	q            []int
	maxSize      int
	lock         sync.Mutex
	qWaitingPuts []*chan int
	qWaitingGets []*chan int
}

func NewNonBlockingQ(size int) *NonBlockingQ {
	return &NonBlockingQ{
		maxSize: size,
		q:       make([]int, 0),
	}
}

func (q *NonBlockingQ) EnQueue(val int) chan int {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.maxSize == len(q.q) {
		fmt.Println("Queue is Full")
		ch := make(chan int)
		q.qWaitingPuts = append(q.qWaitingPuts, &ch)
		return ch
	}
	q.q = append(q.q, val)

	if len(q.qWaitingGets) > 0 {
		ch := q.qWaitingGets[0]
		q.qWaitingGets = q.qWaitingGets[1:]
		*ch <- val
	}

	return nil
}

func (q *NonBlockingQ) DeQueue() (int, chan int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.q) != 0 {
		item := q.q[0]
		q.q = q.q[1:]

		if len(q.qWaitingPuts) > 0 {
			ch := q.qWaitingPuts[0]
			q.qWaitingPuts = q.qWaitingPuts[1:]
			*ch <- 1
		}
		return item, nil
	}

	ch := make(chan int)
	q.qWaitingGets = append(q.qWaitingGets, &ch)
	return 0, ch
}

func retryEnque(ch <-chan int, val int, q *NonBlockingQ) {
	<-ch
	fmt.Println("Retry Enque Invoked ")
	newCh := q.EnQueue(val)
	if newCh != nil {
		go retryEnque(newCh, val, q)
	} else {
		fmt.Printf("\nitem %d successfully added on a retry\n", val)
	}

}

func retryDeque(ch <-chan int) {
	item := <-ch
	fmt.Println("Retry Deque Invoked consumed item", item)
}

func Producer(q *NonBlockingQ) {
	item := 1
	for {
		ch := q.EnQueue(item)
		if ch != nil {
			go retryEnque(ch, item, q)
		}
		item++
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
	}
}

func Consumer(q *NonBlockingQ) {
	for {

		item, ch := q.DeQueue()
		if ch != nil {
			go retryDeque(ch)
		} else {
			fmt.Printf("\nConsumer consumed item %d\n", item)
		}

		time.Sleep(time.Second)
	}
}

func main() {
	q := NewNonBlockingQ(2)
	//go Producer(q)
	//go Consumer(q)

	//time.Sleep(15 * time.Second)
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	go func(){
		defer wg.Done()
		Producer(q)
	}
	
	go func(){
		defer wg.Done()
		Consumer(q)
	}
	
	wg.Wait()	
}
```
