package test

import (
	"editor/core"
	"testing"
)

func BenchmarkBufferInsertBack(b *testing.B) {
	buffer := core.NewEmptyBuffer()
	for i := 0; i < b.N; i++ {
		buffer.Insert([]rune{'x'}, 0)
	}
}

func BenchmarkBufferInsertFront(b *testing.B) {
	buffer := core.NewEmptyBuffer()
	for i := 0; i < b.N; i++ {
		buffer.Insert([]rune{'x'}, i)
	}
}
