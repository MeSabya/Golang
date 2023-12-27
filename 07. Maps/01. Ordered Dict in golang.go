package main

import "fmt"

type node struct {
	key   string
	value interface{}
	prev  *node
	next  *node
}

type OrderedDict struct {
	head  *node
	tail  *node
	nodes map[string]*node
}

func NewOrderedDict() *OrderedDict {
	return &OrderedDict{
		nodes: make(map[string]*node),
	}
}

func (o *OrderedDict) Set(key string, value interface{}) {
	if existingNode, exists := o.nodes[key]; exists {
		// Update existing node
		existingNode.value = value
		return
	}

	newNode := &node{
		key:   key,
		value: value,
	}

	if o.head == nil {
		// First element
		o.head = newNode
		o.tail = newNode
	} else {
		// Append to the end of the list
		newNode.prev = o.tail
		o.tail.next = newNode
		o.tail = newNode
	}

	o.nodes[key] = newNode
}

func (o *OrderedDict) Get(key string) interface{} {
	if node, exists := o.nodes[key]; exists {
		return node.value
	}
	return nil
}

func (o *OrderedDict) Keys() []string {
	keys := make([]string, 0, len(o.nodes))
	current := o.head
	for current != nil {
		keys = append(keys, current.key)
		current = current.next
	}
	return keys
}

func main() {
	// Create a new ordered map
	orderedMap := NewOrderedDict()

	// Set values in order
	orderedMap.Set("one", 1)
	orderedMap.Set("two", 2)
	orderedMap.Set("three", 3)

	// Print keys and values in order
	keys := orderedMap.Keys()
	for _, key := range keys {
		fmt.Printf("%s: %v\n", key, orderedMap.Get(key))
	}
}
