package main

import "testing"

func sum(a, b int) int {
	return a + b
}

//go:noinline
func sumNoInline(a, b int) int {
	return a + b
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum(1, 2)
	}
}

func BenchmarkSumNoInline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumNoInline(1, 2)
	}
}
