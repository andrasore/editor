package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testBuffer struct {
	text []rune
}

func (b *testBuffer) Read(from, to int) []rune {
	return b.text[from:to]
}

func (b *testBuffer) Size() int {
	return len(b.text)
}

func (b *testBuffer) Insert([]rune, int) {
	panic("readonly buffer")
}

func (b *testBuffer) PutChar(rune, int) {
	panic("readonly buffer")
}

func (b *testBuffer) Delete(int, int) {
	panic("readonly buffer")
}

const content = "abcde\n" +
	"abcd\n" +
	"abc\n" +
	"ab"

func TestBufferView_Update(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := defaultBufferView{buffer: &tb}
	bw.Update(0, tb.Size())
	assert.Equal(t, 4, bw.GetLineCount())
	assert.Equal(t, "abcde", string(bw.GetLine(0)))
	assert.Equal(t, "abcd", string(bw.GetLine(1)))
	assert.Equal(t, "abc", string(bw.GetLine(2)))
	assert.Equal(t, "ab", string(bw.GetLine(3)))
}

func TestBufferView_GetCursorIndex(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := defaultBufferView{buffer: &tb}
	bw.Update(0, tb.Size())
	assert.Equal(t, 2, string(bw.GetLine(3)))
}