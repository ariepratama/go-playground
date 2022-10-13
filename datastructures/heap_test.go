package datastructures

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	heap := NewIntHeap()
	heap.Insert(1)
	heap.Insert(2)
	heap.Insert(3)
	heap.Insert(4)
	heap.Insert(5)
	heap.Insert(6)
	heap.Insert(7)
	heap.Insert(8)
	heap.Insert(9)
	heap.Print()

	fmt.Printf("max element is %d\n", heap.GetMax())
}
