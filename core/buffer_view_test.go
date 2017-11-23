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

type coord struct {
	line, char int
}

type expected struct {
	index int
	char  rune
}

func TestBufferView_GetCursorIndex(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := GetBufferView(tb)

	assertCursor := func(expected expected, testInput coord) {
		index, _ := bw.GetCursorIndex(testInput.line, testInput.char)
		assert.Equal(t, expected.index, index)
		assert.Equal(t, expected.char, tb.Read(index, index+1)[0])
	}

	//first index should be zero
	assertCursor(expected{0, 'a'}, coord{0, 0})

	//index in first line is equal to char number
	assertCursor(expected{3, 'd'}, coord{0, 3})

	//index after last char in line should return index of newline
	assertCursor(expected{5, '\n'}, coord{0, 5})

	//all lines should behave in a similar manner
	assertCursor(expected{6, 'a'}, coord{1, 0})
	assertCursor(expected{9, 'd'}, coord{1, 3})
	assertCursor(expected{10, '\n'}, coord{1, 4})
	assertCursor(expected{13, 'c'}, coord{2, 2})
	//TODO - if newline is the last char
}

func TestBufferView_GetCursorIndex_EmptyBuffer(t *testing.T) {
	tb := testBuffer{}
	bw := defaultBufferView{buffer: &tb}
	index, _ := bw.GetCursorIndex(0, 0)
	assert.Equal(t, 0, index) //line 0, char 3: d
}
