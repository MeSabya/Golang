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



