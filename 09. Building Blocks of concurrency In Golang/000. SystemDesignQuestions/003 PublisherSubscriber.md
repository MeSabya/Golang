```golang
package main

import (
	"fmt"
	"sync"
)

type PubSub interface {
	Publish(topic string, message interface{})
	Subscribe(topic string) <-chan interface{}
	Wait()
}

type PubSubImpl struct {
	waitGroup        sync.WaitGroup
	topics           map[string][]chan interface{}
	subscriptionLock sync.Mutex
}

func NewPubSub() *PubSubImpl {
	return &PubSubImpl{
		topics: make(map[string][]chan interface{}),
	}
}

func (ps *PubSubImpl) Publish(topic string, message interface{}) {
	ps.subscriptionLock.Lock()
	defer ps.subscriptionLock.Unlock()

	subscribers := ps.topics[topic]
	for _, subscriber := range subscribers {
		ps.waitGroup.Add(1)
		go func(subscriber chan interface{}) {
			msg := fmt.Sprintf("%s %v", topic, message)
			subscriber <- msg
			ps.waitGroup.Done()
		}(subscriber)
	}

}

func (ps *PubSubImpl) Wait() {
	ps.waitGroup.Wait()
}

func (ps *PubSubImpl) Subscribe(topic string) <-chan interface{} {
	ps.subscriptionLock.Lock()
	defer ps.subscriptionLock.Unlock()

	subscriber := make(chan interface{})
	ps.topics[topic] = append(ps.topics[topic], subscriber)

	return subscriber
}

var pubsub PubSub

func main() {
	pubsub = NewPubSub()

	subscriber1 := pubsub.Subscribe("topic1")
	subscriber2 := pubsub.Subscribe("topic2")
	subscriber3 := pubsub.Subscribe("topic3")
	subscriber4 := pubsub.Subscribe("topic3")
	subscriber5 := pubsub.Subscribe("topic3")

	// Publish a message to the topics
	pubsub.Publish("topic1", "Hello, subscribers number one!")
	pubsub.Publish("topic1", "Bye, subscribers number one!")
	pubsub.Publish("topic2", "Hello, subscribers number two!")
	pubsub.Publish("topic2", "How are you? subscribers number two!")
	pubsub.Publish("topic2", "Bye, subscribers number two!")
	pubsub.Publish("topic3", "Hello, subscribers number three!")
	pubsub.Publish("topic3", "How are you? subscribers number three!")
	pubsub.Publish("topic3", "Bye, subscribers number three!")

	// Receive messages from different topics
	//Wait on all the subscription channels to receive messages
	go func() {
		for {
			select {
			case message := <-subscriber1:
				fmt.Println("subcriber 1", message)
			case message := <-subscriber2:
				fmt.Println("subcriber 2", message)
			case message := <-subscriber3:
				fmt.Println("subcriber 3", message)
			case message := <-subscriber4:
				fmt.Println("subcriber 4", message)
			case message := <-subscriber5:
				fmt.Println("subcriber 5", message)
			}

		}
	}()

	// Wait for all messages to be received by subscribers
	pubsub.Wait()

}
```
## Why the publish function needs to spawn goroutine .. what is the need?

Reasons for Spawning Goroutines in Publish

### 1. Non-Blocking Behavior
Publishing a message to multiple subscribers can potentially block if any of the subscriber channels are not ready to receive the message. 
By spawning a goroutine for each message delivery, you ensure that the Publish function does not block waiting for individual subscribers to receive the message. 
This is crucial for maintaining the responsiveness and throughput of the publisher.

### 2. Concurrency
Using goroutines allows multiple subscribers to process messages concurrently. This parallel processing can lead to better performance and more efficient use of 
system resources, especially when the number of subscribers is large or when processing messages involves time-consuming tasks.


