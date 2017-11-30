package core

import (
	"fmt"
)

const (
	screenWidth  = 19
	screenHeight = 5
)

type testScreen struct {
	content [screenHeight][screenWidth]rune
	cursorX int
	cursorY int
}

func (s *testScreen) Print() {
	//s.content[s.cursorX][s.cursorY] = 'C'
	for _, line := range s.content {
		fmt.Println(string(line[:]))
	}
}

func (s *testScreen) SetCell(x, y int, c rune, fg, bg int) {
	s.content[y][x] = c
}

func (s *testScreen) SetCursor(x, y int) {
	s.cursorX = x
	s.cursorY = y
}

func (s *testScreen) Size() (width int, height int) {
	return screenWidth, screenHeight
}

func (s *testScreen) Clear() {
	s.content = [screenHeight][screenWidth]rune{}
}

func (s *testScreen) Flush() {
}

func getEditor(s *testScreen) Editor {
	buffer := NewEmptyBuffer()
	bufferView := NewBufferView(buffer)
	return Editor{
		Screen:     s,
		Buffer:     buffer,
		BufferView: bufferView,
	}
}

func ExamplePutChar() {
	s := testScreen{}
	ed := getEditor(&s)
	ed.SendChar('i')
	ed.SendChar('x')
	ed.SendChar(KeyBackspace)
	ed.SendChar('y')
	s.Print()
	// Output: y
}
