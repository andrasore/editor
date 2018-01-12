package test

import (
	"bufio"
	"editor/core"
	"strings"
	"testing"
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

func ContentToString(content [screenHeight][screenWidth]rune) string {
	var s []string
	for _, line := range content {
		s = append(s, string(line[:]))
	}
	return strings.Join(s, "\n")
}

func (s *testScreen) Check(t *testing.T, expected string) {
	scanner := bufio.NewScanner(strings.NewReader(expected))
	line := 0
	for scanner.Scan() {
		for i, expectedChar := range scanner.Text() {
			contentChar := s.content[line][i]
			if expectedChar != contentChar {
				contentString := ContentToString(s.content)
				t.Fatalf(
					"\nexpected:\n%v\ngot:\n%v",
					expected,
					contentString[0:len(expected)],
				)
			}
		}
		line++
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

func getEditor(s *testScreen) core.Editor {
	buffer := core.NewEmptyBuffer()
	bufferView := core.NewBufferView(buffer)
	return core.Editor{
		Screen:     s,
		Buffer:     buffer,
		BufferView: bufferView,
	}
}

func TestPutChar(t *testing.T) {
	expected := "y"
	s := testScreen{}
	ed := getEditor(&s)
	ed.SendChar('i')
	ed.SendChar('x')
	ed.SendChar(core.KeyBackspace)
	ed.SendChar('y')
	s.Check(t, expected)
}
