package datastructures

import "fmt"

// Heap implementation of binary heap
// https://www.geeksforgeeks.org/max-heap-in-java/
type Heap[T int | uint] struct {
	Values []T
}

func NewIntHeap() *Heap[int] {
	return &Heap[int]{}
}

func (heap *Heap[T]) GetMax() T {
	return heap.Values[0]
}

func (heap *Heap[T]) Insert(element T) {
	fmt.Printf("inserting... %d\n", element)
	heap.Values = append(heap.Values, element)
	current := len(heap.Values) - 1
	for heap.Values[current] > heap.Values[heap.parent(current)] {
		heap.swap(current, heap.parent(current))
		current = heap.parent(current)
	}
}

func (heap *Heap[T]) Print() {
	for i := 0; i < len(heap.Values)/2; i++ {
		fmt.Printf("Parent Node: %d ", heap.Values[i])
		if heap.leftChild(i) < len(heap.Values) {
			fmt.Printf("Left child node: %d ", heap.Values[heap.leftChild(i)])
		}

		if heap.rightChild(i) < len(heap.Values) {
			fmt.Printf("Right child node: %d ", heap.Values[heap.rightChild(i)])
		}
		fmt.Println()
	}
}

// parent of a pos is always the element seated on half of the position
func (heap *Heap[T]) parent(pos int) int {
	return (pos - 1) / 2
}

// leftChild return the left child position of the current position
// [A,B,C,D,E,F] -> left child of B is 2 * pos + 1 => D
func (heap *Heap[T]) leftChild(pos int) int {
	return (2 * pos) + 1
}

// rightChild return the left child position of the current position
// [A,B,C,D,E,F] -> left child of B is 2 * pos + 2 => E
func (heap *Heap[T]) rightChild(pos int) int {
	return (2 * pos) + 2
}

// swap 2 element position
func (heap *Heap[T]) swap(arrPos1 int, arrPos2 int) {
	temp := heap.Values[arrPos1]
	heap.Values[arrPos1] = heap.Values[arrPos2]
	heap.Values[arrPos2] = temp
}
