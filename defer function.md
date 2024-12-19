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
