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


## Defer types: Heap-allocated, Stack-allocated and Open-coded 
### Heap-Allocated Defer
- Description: When a defer statement is complex, such as when it's used in a loop or involves closures, the deferred call is allocated on the heap.
- Use Case: Happens when the number of defer calls or their lifetimes cannot be determined at compile time.

#### Example1: Using defer inside the loop
```go
func heapAllocatedDefer() {
	for i := 0; i < 10; i++ {
		defer fmt.Println("Deferred in loop:", i) // Cannot be optimized due to dynamic loop behavior.
	}
}
```
#### Example2: Using a Closure in defer
```go
func deferWithClosure() {
	message := "Hello"
	defer func() {
		fmt.Println("Deferred with closure:", message)
	}()
	message = "World"
}
```
#### Dynamic defer Calls
When the arguments or function being deferred are determined dynamically, Go must allocate the deferred call on the heap.

Example:

```go
func deferDynamicCall(flag bool) {
	var message string
	if flag {
		message = "Dynamic defer: true"
	} else {
		message = "Dynamic defer: false"
	}

	defer fmt.Println(message)
}
```
#### defer Inside a Nested Function
When a defer is used in a nested function, the compiler often allocates it on the heap.

Example:

```go
func deferInNestedFunction() {
	innerFunc := func(msg string) {
		defer fmt.Println("Deferred in nested function:", msg)
	}
	innerFunc("Hello, nested!")
}
```
#### Defer in a Goroutine
When defer is used inside a goroutine, the deferred call is heap-allocated because the execution context spans a separate goroutine.

Example:

```go
func deferInGoroutine() {
	go func() {
		defer fmt.Println("Deferred in goroutine")
	}()
}
```
![image](https://github.com/user-attachments/assets/1ac94151-5bc4-4a77-8081-f0014f327b76)

### Stack-Allocated Defer
Description: If the defer call is simple and the compiler can guarantee it will be executed in the same function context, it can be allocated on the stack.
Use Case: Applies to straightforward, non-looped defer calls.
Performance: More efficient because no heap allocation or garbage collection is involved.
Example:

```go
func stackAllocatedDefer() {
	defer fmt.Println("Simple defer")
	fmt.Println("Doing something")
}
```
Explanation:

The compiler optimizes the defer since it knows exactly when and where it will execute.
The defer call is managed using the stack, avoiding heap overhead.

### Open-Coded Defer
Description: In Go 1.20 and later, simple defer calls in certain functions may be inlined or "open-coded." This optimization avoids any runtime defer mechanism by directly inserting the deferred call at the end of the function.
Use Case: Applies to trivial defer calls where there is no looping or closure involved, and the function body is small.
Performance: Fastest because it eliminates the runtime cost of managing a defer call entirely.
Example:

```go
func openCodedDefer() {
	defer fmt.Println("Open-coded defer")
	fmt.Println("Hello, world!")
}
```
Explanation:

- The compiler inserts the fmt.Println("Open-coded defer") call directly before the function return.
- This avoids both stack and heap allocation, leading to zero runtime overhead for defer.

