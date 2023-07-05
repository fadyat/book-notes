package main

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestAsFunctionArgument(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	func(m map[string]int) {
		for i := 0; i < 100; i++ {
			m[strconv.Itoa(i)] = i
		}
	}(m)

	require.NotEqual(t, m, map[string]int{"one": 1, "two": 2})
}

func TestCantTakeAddress(t *testing.T) {
	// m := map[string]int{"one": 1, "two": 2}

	// This will not compile, because evaluation may happen, and address will change
	// a := &m["one"]
	// _ = a

	t.Skip("This test will not compile")
}

func TestMapPreallocation(t *testing.T) {
	m := make(map[string]int, 100)
	require.Equal(t, 0, len(m))

	// no need to preallocate, because map already preallocated
	for i := 0; i < 100; i++ {
		m[strconv.Itoa(i)] = i
	}

	require.Equal(t, 100, len(m))
}
