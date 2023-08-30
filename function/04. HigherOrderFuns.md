## Function used as a value 

```golang
func inc1(x int) int { return x+1 }
f1 := inc1 // f1 := func (x int) int { return x+1 }
```

## Function used as a parameter

```golang
package main
import (
    "fmt"
)

func main() {
    callback(1, Add) // function passed as a parameter
}

func Add(a, b int) {
    fmt.Printf("The sum of %d and %d is: %d\n", a, b, a + b)
}

func callback(y int, f func(int, int)) {
    f(y, 2) // this becomes Add(1, 2)
}
```

## Function used as a filter 

```golang
package main
import "fmt"

type flt func(int) bool
    // isOdd takes an int slice and returns a bool set to true if the
    // int parameter is odd, or false if not.
    // isOdd is of type func(int) bool which is what flt is declared to be.

func isOdd(n int) bool {
        if n % 2 == 0 {
            return false
        }
        return true
    }

    // Same comment for isEven
func isEven(n int) bool {
    if n % 2 == 0 {
        return true
    }
    return false
}

func filter(sl[] int, f flt)[] int {
    var res[] int
    for _, val := range sl {
        if f(val) {
            res = append(res, val)
        }
    }
    return res
}

func main() {
    slice := [] int {1, 2, 3, 4, 5, 7}
    fmt.Println("slice = ", slice)
    odd := filter(slice, isOdd)
    fmt.Println("Odd elements of slice are: ", odd)
    even := filter(slice, isEven)
    fmt.Println("Even elements of slice are: ", even)
}
```
