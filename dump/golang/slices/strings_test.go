package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCorrectString(t *testing.T) {
	bytes := []byte("Hello World")
	str := string(bytes)

	assert.Equal(t, "Hello World", str)
}

// When we calling string() on a slice of bytes, it will make memory allocation
// So the string will be a copy of the slice of bytes
func TestBytesToString(t *testing.T) {
	bytes := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	str := string(bytes)

	assert.Equal(t, "Hello", str)

	bytes[0] = 0x41
	assert.Equal(t, "Hello", str)

	str += " World"
	assert.Equal(t, "Hello World", str)
}

// Golang saving string in UTF-8 encoding, the real symbol can be more than 1 byte
//
// When using range on a string, it will return the index of the symbol and the symbol itself
// Slicing on string works by byte index, not symbol index
func TestUTF8Encoding(t *testing.T) {
	utf8 := "Hello 世界"
	assert.NotEqual(t, 8, len(utf8))

	i := 0
	for range utf8 {
		i++
	}

	assert.Equal(t, 8, i)

	cut := utf8[:8]
	assert.NotEqual(t, "Hello 世界", cut)
}

// Runes are the actual symbols in the string, not the bytes
// Indexing on runes works by symbol index, not byte index
func TestRune(t *testing.T) {
	utf8 := "Hello 世界"

	runes := []rune(utf8)
	assert.Equal(t, 8, len(runes))
	assert.Equal(t, '世', runes[6])

	runes = append(runes, []rune{'э', 'й'}...)
	assert.Equal(t, "Hello 世界эй", string(runes))
}
