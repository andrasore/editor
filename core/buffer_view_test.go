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
	assert.Equal(t, 5, bw.LineCount())
	assert.Equal(t, "abcde", string(bw.Line(0)))
	assert.Equal(t, "abcd", string(bw.Line(1)))
	assert.Equal(t, "abc", string(bw.Line(2)))
	assert.Equal(t, "ab", string(bw.Line(3)))
	assert.Equal(t, "", string(bw.Line(4)))
}

func TestBufferView_GetLine_WithNoNewline(t *testing.T) {
	tb := testBuffer{[]rune("abcd")}
	bw := NewBufferView(tb)
	assert.Equal(t, "abcd", string(bw.Line(0)))
}

type coord struct {
	line, char int
}

func TestBufferView_IndexOf(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := NewBufferView(tb)

	assertIndex := func(expected int, testInput coord) {
		t.Helper()
		index := bw.IndexOf(testInput.line, testInput.char)
		assert.Equal(t, expected, index)
	}

	//index in first line is equal to char number
	assertIndex(0, coord{0, 0})
	assertIndex(3, coord{0, 3})

	//index after last char in line should return index of newline
	assertIndex(5, coord{0, 5})

	//all lines should behave in a similar manner
	assertIndex(6, coord{1, 0})
	assertIndex(9, coord{1, 3})
	assertIndex(10, coord{1, 4})
	assertIndex(13, coord{2, 2})
}

func TestBufferView_IndexOf_EmptyBuffer(t *testing.T) {
	tb := testBuffer{}
	bw := defaultBufferView{buffer: &tb}
	position := bw.IndexOf(0, 0)
	assert.Equal(t, 0, position)
}

func TestBufferView_IndexOf_Overindexing(t *testing.T) {
	tb := testBuffer{[]rune(content)}
	bw := NewBufferView(tb)
	assert.Panics(t, func() { bw.IndexOf(0, 6) })
	assert.Panics(t, func() { bw.IndexOf(1, 5) })
	assert.Panics(t, func() { bw.IndexOf(2, 4) })
}

const emptyLinesContent = "ab\n" + "\n" + "\n"

func TestBufferView_IndexOf_EmptyLines(t *testing.T) {
	tb := testBuffer{[]rune(emptyLinesContent)}
	bw := NewBufferView(tb)
	assert.Equal(t, 3, bw.IndexOf(1, 0))
	assert.Equal(t, 4, bw.IndexOf(2, 0))
}

func TestBufferView_PositionOf(t *testing.T) {
	const content = "abc\n" +
		"ab\n"

	tb := testBuffer{[]rune(content)}
	bw := NewBufferView(tb)

	testCases := []struct {
		line  int
		char  int
		index int
	}{
		{line: 0, char: 0, index: 0},
		{line: 0, char: 0, index: 3},
		{line: 1, char: 0, index: 4},
		{line: 2, char: 0, index: 7},
	} //we should be able to get position after last newline

	for _, testCase := range testCases {
		resultLine, resultChar := bw.PositionOf(testCase.index)
		if resultLine != testCase.line || resultChar != testCase.char {
			t.Errorf(
				"Expected line: %v, char: %v Got: %v, %v",
				testCase.line,
				testCase.char,
				resultLine,
				resultChar,
			)
		}
	}
}
