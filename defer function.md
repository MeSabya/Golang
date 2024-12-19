In Go, defer is a keyword used to delay the execution of a function until the surrounding function finishes.

```go
func main() {
  defer fmt.Println("hello")
  fmt.Println("world")
}

// Output:
// world
// hello
```
## Key Points About defer and panic:
### Defer Registration:
- The defer statement registers the deferred function at the moment it is encountered in the code.
- If a panic occurs before the defer is registered, that deferred function will not execute.

### Defer Execution During Panic:
- When a panic occurs, Go starts unwinding the stack of function calls.
- During this unwinding process, any deferred functions in the affected stack frames are executed in the reverse order of their registration.

### Ensuring Defer Execution:
- To ensure that a defer runs even during a panic, it must be placed before any code that might cause a panic.

### Example to illustrate the above

```go
package main

import "fmt"

func examplePanic() {
	defer fmt.Println("Deferred function executed.")
	fmt.Println("About to panic!")
	panic("Something went wrong!")
}

func noDeferPanic() {
	fmt.Println("Panic without defer!")
	panic("No defer registered before this.")
}

func main() {
	fmt.Println("=== Running example with defer ===")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from:", r)
			}
		}()
		examplePanic()
	}()

	fmt.Println("\n=== Running example without defer ===")
	func() {
	    noDeferPanic()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from:", r)
			}
		}()
		
	}()
}
```

👉 **The recover function is meant to catch a panic, but it has to be called within a deferred function to work properly.**

### Tricky Question1: Does this code will work?
```go
func myRecover() {
  if r := recover(); r != nil {
    fmt.Println("Recovered:", r)
  }
}

func main() {
  defer func() {
    myRecover()
    // ...
  }()

  panic("This is a panic")
}
```
the code above won’t work as you might expect. That’s because recover isn’t called directly from a deferred function but from a nested function.

### Tricky Question 2: Does this code will work?

```go
func main() {
  defer func() {
    if r := recover(); r != nil {
      fmt.Println("Recovered:", r)
    }
  }()

  go panic("This is a panic")

  time.Sleep(1 * time.Second) // Wait for the goroutine to finish
}
```
one goroutine could not intervene in another to handle the panic since each goroutine has its own stack.


## Defers are stacked #
When you use multiple defer statements in a function, they are executed in a ‘stack’ order, meaning the last deferred function is executed first.

```go
func main() {
  defer fmt.Println(1)
  defer fmt.Println(2)
  defer fmt.Println(3)
}

// Output:
// 3
// 2
// 1
```
Every time you call a defer statement, you’re adding that function to the top of the current goroutine’s linked list, like this:

![image](https://github.com/user-attachments/assets/5555308a-1cb7-4815-ab1f-a21b456f9f13)

```go
func B() {
  defer fmt.Println(1)
  defer fmt.Println(2)
  A()
}

func A() {
  defer fmt.Println(3)
  defer fmt.Println(4)
}
```
![image](https://github.com/user-attachments/assets/e4cb0a63-8950-43a5-b433-0244a479eb34)





