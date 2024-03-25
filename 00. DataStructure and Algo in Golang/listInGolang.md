
```golang
package main

import (
	"container/list"
	"fmt"
)

type Item struct {
	key int
	value interface{}
}


type LRUCache struct {
	Capacity int
	itemsList *list.List
	itemsMap map[int]*list.Element	
}

func Constructor(capacity int) LRUCache {
	lru := LRUCache{
		Capacity: capacity,
		itemsList: list.New(),
		itemsMap: make(map[int]*list.Element)
	}

	return lru
}

func (this *LRUCache) Get(key int) int {
	//Case1: Get the node from the map.In any case we need to move the node to the front.
	if node, found := this.itemsMap[key]; ok {
		this.itemsList.MoveToFront(node)
		return node.Value.(Item).value.(int)
	}
	
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if node, found := this.itemsMap[key]; found {
		node.Value = Item{key, value}
		this.itemsList.MoveToFront(node)	
		return
	}
	
	if this.Capacity == len(this.itemsList) {
		delete(this.itemsMap, this.itemsList.Back().Value.(Item).key)
		this.itemsList.Remove(this.itemsList.Back())
	}
	
	node = this.itemsList.PushFront(Item{key, value})
	this.itemsMap[key] = node
	
}
```

## What is list.Element here ? What are methods/attributes  available in itemsList.

In Go, list.Element represents an element in a doubly linked list. It's a type defined in the standard library's container/list package. Each list.Element contains 
a Value field that holds the value of the element and Next and Prev fields that point to the next and previous elements in the list, respectively.

***Here are the key methods and attributes available for list.List:***

- PushFront: Adds a new element with the specified value to the front of the list.
- PushBack: Adds a new element with the specified value to the back of the list.
- Remove: Removes the specified element from the list. This operation takes constant time.
- MoveToFront: Moves the specified element to the front of the list, making it the most recently used.
- Len: Returns the number of elements in the list.
- Front: Returns the first element of the list.
- Back: Returns the last element of the list.
- Init: Initializes or clears the list.
- InsertAfter: Inserts a new element with the specified value after a given element.
- InsertBefore: Inserts a new element with the specified value before a given element.
- MoveAfter: Moves a specified element to after another specified element.
- MoveBefore: Moves a specified element to before another specified element.



