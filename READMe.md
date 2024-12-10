<details>
<summary>How a golang code project structure should look like</summary>

## How a golang code project structure should look like

In Go, the recommended location for your Go source code is within the Go workspace. The Go workspace typically has the following directory structure:

```lua
GOPATH/
|-- bin/
|-- pkg/
|-- src/
    |-- github.com/
        |-- yourusername/
            |-- yourproject/
                |-- main.go
```
Here:

- GOPATH is an environment variable that points to the root of your workspace.
- bin/ contains the executable binaries.
- pkg/ contains package objects.
- src/ contains your source code.

Each project you work on should have its own folder under src/. In the example above, yourusername represents your GitHub username, and yourproject is the name of your Go project.

However, starting from Go 1.11, you have the option to work outside the traditional workspace by using Go modules. With Go modules, you can place your code in any directory, and Go will manage dependencies for you.

Here is an example of a project structure using Go modules:

```lua
myproject/
|-- go.mod
|-- go.sum
|-- main.go
|-- greetings/
    |-- greetings.go
```

In this structure:

- go.mod and go.sum are files created and managed by Go modules.
- main.go is your main program.
- greetings/ is a package that you might create.

To create a Go module, you can run the following command inside your project directory:

```bash
go mod init myproject
```

This will initialize a Go module for your project.

Remember, Go is flexible, and you have the freedom to organize your code the way that makes sense for your project. The Go module approach provides more flexibility in terms of project organization and dependency management.

### How can you tell Go to import a package from a different location?

In Go, you can use the import statement to import packages from different locations. By default, Go imports packages from the Go module specified in the go.mod file or from the standard library. However, you can specify a different import path for a package if it is hosted in a different location (e.g., a different repository or a custom server).

To import a package from a different location, you need to provide the full import path in your source code. The import path is a unique identifier for a package that includes the module name and the path within the module where the package is located.

Here's the general syntax:

```go
import "module/path/package"
``
Here's an example:

```go
// Importing a package from a different location
import "github.com/example/mylibrary/mypackage"
```
In this example, github.com/example/mylibrary is the module path, and mypackage is the path within the module where the package is located.

If the package is not part of a Go module, you can use the full URL of the repository:

```go
// Importing a package from a GitHub repository
import "github.com/example/mylibrary/mypackage"
```
</details>


<details>
    <summary>What do you need for two functions to be the same type?</summary>
    
In Go, for two functions to be considered the same type, they must have the same parameter types, the same return types, and the same names for corresponding parameters (if named parameters are used). The function signatures, which include the parameter and return types, need to match exactly.

Here's an example:
```golang
package main

import "fmt"

// Function1 has the same type as Function2
func Function1(a int, b string) {
    fmt.Println("Function1:", a, b)
}

func Function2(x int, y string) {
    fmt.Println("Function2:", x, y)
}

func main() {
    // Both function variables have the same type
    var f1 func(int, string) = Function1
    var f2 func(int, string) = Function2

    f1(42, "hello")
    f2(42, "world")
}
```
</details>

<details>
    <summary>From where is the variable myVar accessible if it is declared outside of any functions in a file in package myPackage located inside module myModule?</summary>
    
In Go, when a variable is declared outside of any functions within a file in a package, it becomes a package-level variable. The accessibility of a package-level variable depends on its identifier's casing (uppercase or lowercase).

## Here are the rules:

### Uppercase (exported) identifier:

If the variable name starts with an uppercase letter (e.g., MyVar), it is considered an exported identifier and is accessible from outside the package.
```go
// mypackage/mypackage.go
package mypackage

var MyVar int = 42
```

```go
// main.go
package main

import "mypackage"

func main() {
    value := mypackage.MyVar
    // You can access MyVar from outside the package because it is uppercase
    println(value)
}
```

### Lowercase identifier:

If the variable name starts with a lowercase letter (e.g., myVar), it is considered unexported and is only accessible within the same package.

```go
// mypackage/mypackage.go
package mypackage

var myVar int = 42
```

```go
// main.go
package main

import "mypackage"

func main() {
    // This would result in a compilation error
    value := mypackage.myVar
    println(value)
}
```
So, the accessibility of MyVar or myVar depends on whether the first letter of the identifier is uppercase (exported) or lowercase (unexported) and whether it is being accessed from within or outside the package.
</details>


<details>
    <summary> Fix the code-1 </summary>
    
```go
type Point struct {
  x int
  y int
}
 
func main() {
  data := []byte(`{"x":1, "y": 2}`)
  var p Point
  if err := json.Unmarshal(data, &p); err != nil {
    fmt.Println("error: ", err)
  } else {
    fmt.Println(p)
  }
}

This code printed {0, 0}. How can you fix it?
```

The issue with the provided code is related to the visibility of the fields in the Point struct. In Go, fields with a lowercase initial letter (e.g., x and y) are unexported and not accessible outside the package where the struct is defined.

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	data := []byte(`{"x":1, "y": 2}`)
	var p Point

	if err := json.Unmarshal(data, &p); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(p)
	}
}
```
    
</details>

<details>
	<summary>Fix the code-2</summary>
	What will be printed in this code?

```go
var stocks map[string]float64 // stock -> price
price := stocks["MSFT"]
fmt.Printf("%f\n", price)
```
       
The provided code will result in a runtime panic. This is because the stocks map is declared but not initialized before attempting to access the value associated with the key "MSFT".

In Go, a map is a reference type, and it must be initialized before use. The zero value of a map is nil, and attempting to access a key in a nil map results in a runtime panic.

To fix this issue, you need to initialize the stocks map before attempting to access its values. Here's an example:

```go
package main

import "fmt"

func main() {
    var stocks map[string]float64 // stock -> price

    // Initialize the map before using it
    stocks = make(map[string]float64)

    // Accessing the value for the key "MSFT"
    price := stocks["MSFT"]
    fmt.Printf("%f\n", price)
}
```
</details>

<details>
	<summary>Fix the code-3</summary>

Given the definition of worker below, what is the right syntax to start a goroutine that will call worker and send the result to a channel named ch?

func worker(m Message) Result

```go
package main

import "fmt"

type Message struct {
	Text string
}

type Result struct {
	ResultText string
}

func worker(m Message) Result {
	// Some processing...
	return Result{ResultText: "Processed: " + m.Text}
}

func main() {
	// Create a channel
	ch := make(chan Result)

	// Create a Message
	message := Message{Text: "Hello, World!"}

	// Start a goroutine to call worker and send the result to the channel
	go func() {
		result := worker(message)
		ch <- result
		close(ch) // Close the channel when done sending
	}()

	// Retrieve the result from the channel
	result := <-ch
	fmt.Println(result.ResultText)
}
```
</details>

### what is meaning of this []int(nil) 
<details>
	<Summary>Answer</Summary>
nil is the zero value for reference types in Go (pointers, slices, maps, channels, and interfaces).
[]int(nil) explicitly converts nil to a slice of type []int. This ensures that the type of the slice is clear, even though the slice itself is nil.

## Why Use []int(nil)?

- To initialize a nil slice explicitly.
- To reset a slice to its nil value.
- To create a slice with a specific type without allocating memory.

## Usage Examples:
### 1. Resetting a Slice:

```go
var slice []int = []int{1, 2, 3}
fmt.Println(slice) // Output: [1 2 3]

// Reset the slice to nil
slice = []int(nil)
fmt.Println(slice) // Output: []
fmt.Println(slice == nil) // Output: true
```
Here, []int(nil) explicitly sets slice to a nil slice of type []int.

### 2. Deep Copy of a Slice:
The expression append([]int(nil), originalSlice...) is a common idiom for creating a deep copy of a slice:

[]int(nil) creates a new empty slice of type []int without any underlying array.
append([]int(nil), originalSlice...) copies all elements of originalSlice into a new slice.
Example:

```go
original := []int{1, 2, 3}

// Create a deep copy
copy := append([]int(nil), original...)
copy[0] = 99

fmt.Println("Original:", original) // Output: Original: [1 2 3]
fmt.Println("Copy:", copy)         // Output: Copy: [99 2 3]
```
Here, modifying copy does not affect original because append creates a new underlying array for the new slice.

![image](https://github.com/user-attachments/assets/e3606c8a-662a-4dc1-bf3b-1f264b81d3e9)

</details>

### Deadlock Example
<details>
	<summary>Deadlock Example</summary>
	
	```go
	package main

import (
	"fmt"
	"sync"
)

func main() {
	tasks := make(chan int)
	var wg sync.WaitGroup

	// Start a consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range tasks {
			fmt.Println("Processing task:", task)
		}
		fmt.Println("Consumer done")
	}()

	// Producer sends tasks
	tasks <- 1
	tasks <- 2
	// Producer waits for the consumer to finish
	wg.Wait() // Deadlock! Channel is not closed, consumer waits forever.
	close(tasks)
}
```
This is classic chicken egg problem , where producer waits for consumer to finish their consumption.
While consumers are waiting on tasks to consume more, since channel is not closed.
</details>


