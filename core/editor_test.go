package core

import (
	"fmt"
	"testing"
)

const (
	screenWidth  = 10
	screenHeight = 5
)

type testScreen struct {
	content [screenWidth][screenHeight]rune
	cursorX int
	cursorY int
}

func (s *testScreen) Print() {
	s.content[s.cursorX][s.cursorY] = 'C'
	for _, line := range s.content {
		fmt.Println(string(line[:]))
	}
}

func (s *testScreen) SetCell(x, y int, c rune, fg, bg int) {
	s.content[x][y] = c
}

func (s *testScreen) SetCursor(x, y int) {
	s.cursorX = x
	s.cursorY = y
}

func (s *testScreen) Size() (width int, height int) {
	return screenWidth, screenHeight
}

func (s *testScreen) Clear() {
	s.content = [screenWidth][screenHeight]rune{}
}

func (s *testScreen) Flush() {
}

func TestEditor(t *testing.T) {
}
