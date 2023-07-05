package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Slicing from slice doesn't copy the data
// It just creates a new slice that points to the same underlying array
// The new slice has a different length and capacity
// * Length = number of elements in the new slice
// * Capacity = number of elements in the underlying array starting from the first element in the new slice
func TestByIndex(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	partialSlice := slice[1:3]

	assert.Equal(t, []int{2, 3}, partialSlice)
	assert.Equal(t, 2, len(partialSlice))
	assert.Equal(t, 4, cap(partialSlice))

	partialSlice[0] = 999
	assert.Equal(t, []int{999, 3}, partialSlice)
	assert.Equal(t, []int{1, 999, 3, 4, 5}, slice)
	assert.Same(t, &slice[1], &partialSlice[0], "slices share the same underlying array")
}

// Because the capacity of the partialSlice is more than the length
// When appending to the partialSlice, it won't create a new array
// It modifies the underlying array, which is also shared by the original slice
// So the original slice is also modified
//
// When we reach the capacity of the partialSlice, it will create a new array
// And the original slice won't be modified
func TestAppend(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	partialSlice := slice[1:3]

	partialSlice = append(partialSlice, 6)
	assert.Equal(t, []int{2, 3, 6}, partialSlice)
	assert.Equal(t, []int{1, 2, 3, 6, 5}, slice)
	assert.Same(t, &slice[3], &partialSlice[2], "slices share the same underlying array")

	partialSlice = append(partialSlice, 7)
	assert.Equal(t, []int{2, 3, 6, 7}, partialSlice)
	assert.Equal(t, []int{1, 2, 3, 6, 7}, slice)
	assert.Same(t, &slice[4], &partialSlice[3], "slices share the same underlying array")

	partialSlice = append(partialSlice, 8)
	assert.Equal(t, []int{2, 3, 6, 7, 8}, partialSlice)
	assert.Equal(t, []int{1, 2, 3, 6, 7}, slice)
	assert.NotSame(t, &slice[4], &partialSlice[3], "slices don't share the same underlying array")
}

// When we create a slice from an array, it's still using the same underlying array
// So when we modify the slice, it will also modify the array
//
// When reaching the capacity of the slice, works the same as in the previous test
func TestCreateFromArray(t *testing.T) {
	slice := [5]int{1, 2, 3, 4, 5}
	copySlice := slice[:]

	assert.Equal(t, [5]int{1, 2, 3, 4, 5}, slice)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, copySlice)
	assert.Same(t, &slice[0], &copySlice[0], "slices share the same underlying array")

	copySlice[0] = 999
	assert.Equal(t, [5]int{999, 2, 3, 4, 5}, slice)
	assert.Equal(t, []int{999, 2, 3, 4, 5}, copySlice)

	copySlice = append(copySlice, 6)
	assert.Equal(t, [5]int{999, 2, 3, 4, 5}, slice)
	assert.Equal(t, []int{999, 2, 3, 4, 5, 6}, copySlice)
	assert.NotSame(t, &slice[0], &copySlice[0], "slices don't share the same underlying array")
}

// When we don't want to share the underlying array, we can use the copy function
// It creates a new array with the same length, capacity and data
// So when we modify the new slice, it won't modify the original slice
func TestCopy(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	copySlice := make([]int, len(slice))
	copy(copySlice, slice)

	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice)

	copySlice[0] = 999
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice)
	assert.Equal(t, []int{999, 2, 3, 4, 5}, copySlice)

	copySlice = append(copySlice, 6)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice)
	assert.Equal(t, []int{999, 2, 3, 4, 5, 6}, copySlice)
}

// When we pass a slice as a function argument, it will pass a copy of the SliceHeader, not the slice itself
//
//	type SliceHeader struct {
//	    Data uintptr
//	    Len  int
//	    Cap  int
//	}
//
// Modification by the index will modify the original slice
// Modification by to append won't update the length of the original slice and we don't see the new element
// When we reach a capacity of the slice, it will create a new array
func TestAsFuncArgument(t *testing.T) {
	slice := make([]int, 1, 3)

	func(slice []int) {
		slice[0] = 999
		assert.Equal(t, []int{999}, slice)
	}(slice)
	assert.Equal(t, []int{999}, slice)

	func(slice []int) {
		slice = append(slice, 111)
		assert.Equal(t, []int{999, 111}, slice)
	}(slice)

	assert.Equal(t, []int{999}, slice, "what we see with old length")
	assert.Equal(t, []int{999, 111}, slice[:2], "update length of the slice")

	slice = append(slice, 222)
	assert.Equal(t, []int{999, 222}, slice, "length isn't updated, will overwrite the hidden element")
}
