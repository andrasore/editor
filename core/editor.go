package core

import "fmt"

const (
	KeyEsc       rune = 0x1B
	KeyEnter     rune = 0x0D
	KeyBackspace rune = 0x08
)

const (
	DirectionUp = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type Editor struct {
	Screen
	Buffer
	BufferView
	window
	state int
}

func (e *Editor) SendChar(c rune) {
	switch e.state {
	case statusInsert:
		switch c {
		case KeyEsc:
			e.state = statusNormal
		case KeyEnter:
			e.NewLine()
		case KeyBackspace:
			e.DeleteCharBefore()
		default:
			e.PutChar(c)
		}
	case statusNormal:
		switch c {
		case 'i':
			e.state = statusInsert
		case 'h':
			e.MoveCursor(DirectionLeft)
		case 'j':
			e.MoveCursor(DirectionDown)
		case 'k':
			e.MoveCursor(DirectionUp)
		case 'l':
			e.MoveCursor(DirectionRight)
		case 'x':
			e.DeleteChar()
		}
	}
	e.BufferView.Update() //TODO - update as necessary
	e.window.redraw(e.Screen, e.BufferView, e.state, e.Buffer.Size())
}

func (e *Editor) getCursorPosition() int {
	char := e.cursor.char
	line := e.cursor.line
	return e.BufferView.IndexOf(line, char)
}

func (e *Editor) NewLine() {
	cursorPosition := e.getCursorPosition()
	e.Buffer.PutChar('\n', cursorPosition)
	e.cursor.char = 0
	e.cursor.line++
}

func (e *Editor) MoveCursor(direction int) {
	line := e.cursor.line
	char := e.cursor.char
	switch direction {
	case DirectionUp:
		if 0 < line {
			prevLineLength := e.BufferView.LineLength(line - 1)
			e.cursor.char = min(max(prevLineLength-1, 0), char)
			e.cursor.line--
		}
	case DirectionDown:
		if line < e.BufferView.LineCount()-1 {
			nextLineLength := e.BufferView.LineLength(line + 1)
			e.cursor.char = min(max(nextLineLength-1, 0), char)
			e.cursor.line++
		}
	case DirectionLeft:
		if 0 < char {
			e.cursor.char--
		}
	case DirectionRight:
		currentLineLength := e.BufferView.LineLength(line)
		if char < currentLineLength-1 {
			e.cursor.char++
		}
	default:
		panic(fmt.Sprintf("Invalid move direction: %v", direction))
	}
}

func (e *Editor) DeleteChar() {
	if e.LineLength(e.cursor.line) == 0 {
		return
	}

	e.Buffer.DeleteChar(e.getCursorPosition())

	if e.LineLength(e.cursor.line) <= e.cursor.char {
		e.cursor.char = max(e.cursor.char-1, 0)
	}
}

func (e *Editor) DeleteCharBefore() {
	cursorIndex := e.getCursorPosition()
	if cursorIndex == 0 {
		return
	}
	cursorIndex--
	newLine, newChar := e.BufferView.PositionOf(cursorIndex)

	e.Buffer.DeleteChar(cursorIndex)

	e.cursor.line = newLine
	e.cursor.char = newChar
}

func (e *Editor) PutChar(c rune) {
	cursorIndex := e.getCursorPosition()
	e.Buffer.PutChar(c, cursorIndex)
	e.cursor.char++
}
