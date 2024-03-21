
## Problem1: Why this can cause deadlock?

```golang
func main() {
	ch := make(chan int)
	ch <- 1
	fmt.Println(<-ch)
}
```

### Buffered channel Vs Unbuffered channel 
In an unbuffered channel, when a sender sends a value, it will be blocked until a receiver is ready to receive that value. Similarly,
when a receiver attempts to receive a value from an unbuffered channel, it will be blocked until there is a sender ready to send a value.

This synchronous communication ensures that there is a direct handoff between the sender and the receiver. It helps in coordinating the 
execution of goroutines by ensuring that both the sender and the receiver are ready before the communication happens.
However, it can lead to deadlocks if there is no synchronization between the sender and receiver.

In the above example there is only one main goroutine, it is blocked while sending the value. So to avoid this situation either we can use a 
buffered channel or we can use another goroutine while sending.

#### Solution1: Buffered channel

```golang

package main

import "fmt"

func main() {
	ch := make(chan int)
	ch <- 1
	fmt.Println(<-ch)
}
```
#### Solution2: Use another goroutine
```golang
package main

import "fmt"

func main() {
    ch := make(chan int)
    go func() {
        ch <- 1 // Send value on channel
    }()
    fmt.Println(<-ch) // Receive value from channel
}
```
## Problem2: 

