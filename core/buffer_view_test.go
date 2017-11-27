package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testBuffer struct {
	text []rune
}

func (b testBuffer) Read(from, to int) []rune {
	return b.text[from:to]
}

func (b testBuffer) Size() int {
	return len(b.text)
}

func (b testBuffer) Insert([]rune, int) {
	panic("readonly buffer")
}

func (b testBuffer) PutChar(rune, int) {
	panic("readonly buffer")
}

func (b testBuffer) Delete(int, int) {
	panic("readonly buffer")
}

const content = "abcde\n" +
	"abcd\n" +
	"abc\n" +
	"ab\n"

func TestBufferView_Update(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := defaultBufferView{buffer: &tb}
	bw.Update(0, tb.Size())
	assert.Equal(t, 5, bw.GetLineCount())
	assert.Equal(t, "abcde", string(bw.GetLine(0)))
	assert.Equal(t, "abcd", string(bw.GetLine(1)))
	assert.Equal(t, "abc", string(bw.GetLine(2)))
	assert.Equal(t, "ab", string(bw.GetLine(3)))
	assert.Equal(t, "", string(bw.GetLine(4)))
}

func TestBufferView_GetLine_WithNoNewline(t *testing.T) {
	tb := testBuffer{[]rune("abcd")}
	bw := NewBufferView(tb)
	assert.Equal(t, "abcd", string(bw.GetLine(0)))
}

type coord struct {
	line, char int
}

func TestBufferView_GetCursorPosition(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := NewBufferView(tb)

	assertCursor := func(expected int, testInput coord) {
		index := bw.GetCursorPosition(testInput.line, testInput.char)
		assert.Equal(t, expected, index)
	}

	//index in first line is equal to char number
	assertCursor(0, coord{0, 0})
	assertCursor(3, coord{0, 3})

	//index after last char in line should return index of newline
	assertCursor(5, coord{0, 5})

	//all lines should behave in a similar manner
	assertCursor(6, coord{1, 0})
	assertCursor(9, coord{1, 3})
	assertCursor(10, coord{1, 4})
	assertCursor(13, coord{2, 2})
}

func TestBufferView_GetCursorPosition_EmptyBuffer(t *testing.T) {
	tb := testBuffer{}
	bw := defaultBufferView{buffer: &tb}
	position := bw.GetCursorPosition(0, 0)
	assert.Equal(t, 0, position)
}

func TestBufferView_GetCursorPosition_Overindexing(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := NewBufferView(tb)
	assert.Panics(t, func() { bw.GetCursorPosition(0, 6) })
	assert.Panics(t, func() { bw.GetCursorPosition(1, 5) })
	assert.Panics(t, func() { bw.GetCursorPosition(2, 4) })
}
