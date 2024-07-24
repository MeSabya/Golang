package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name    string
	Age     int
	Address string
}

func LearnReflect() {
	//What are methods available in reflect package
	//TypeOf
	//ValueOf
	user := User{Name: "John", Age: 30, Address: "New York"}
	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)

	field := t.Field(1)

	fmt.Println("t is", t)
	fmt.Println("v is", v)
	fmt.Println("Field name:", field.Name)
	fmt.Println("Field type:", field.Type)

}

func FlattenList(L []interface{}) []interface{} {
	var result []interface{}

	for _, item := range L {
		if sublist, ok := item.([]interface{}); ok {
			result = append(result, FlattenList(sublist)...)
		} else {
			result = append(result, item)
		}
	}
	return result
}

/*
In Go, type assertions are used to convert an interface type to a specific type. The expression nestedList.([]interface{}) is 
a type assertion that converts nestedList from the empty interface type interface{} to a slice of interface{} ([]interface{}).

This is necessary because interface{} can hold any value of any type, but we need to explicitly tell the compiler that we 
expect nestedList to be a slice of interface{} so that we can iterate over it. Without this type assertion, the compiler does 
not know the underlying type of nestedList.
*/

func FlattenList2(slice []interface{}) []interface{} {
	var result []interface{}

	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		// Not a slice, return a slice with the single element
		return []interface{}{slice}
	}

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i).Interface()

		switch item := item.(type) {
		case []interface{}:
			result = append(result, FlattenList2(item)...)

		default:
			result = append(result, item)
		}

	}

	return result
}

func main() {
	list := []interface{}{1, 2, []interface{}{3, 4}, []interface{}{5, 6}}
	result := FlattenList2(list)
	fmt.Println(result)
}
