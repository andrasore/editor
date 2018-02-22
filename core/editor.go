package core

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
		}
	}

	e.window.redraw(e.Screen, e.Buffer, e.state, e.Buffer.Size())
}

func (e *Editor) getCursorPosition() int {
	char := e.window.cursor.char
	line := e.window.cursor.line
	return e.BufferView.IndexOf(line, char)
}

func (e *Editor) NewLine() {
	cursorPosition := e.getCursorPosition()
	e.Buffer.PutChar('\n', cursorPosition)
	e.window.cursor.char = 0
	e.window.cursor.line++
	e.BufferView.Update()
}

func (e *Editor) MoveCursor(direction int) {
	line := e.window.cursor.line
	char := e.window.cursor.char

	switch direction {
	case DirectionUp:
		if 0 < line {
			e.window.cursor.line--
			prevLineLength := e.BufferView.LineLength(line - 1)
			if prevLineLength-1 < char {
				e.window.cursor.char = prevLineLength - 1
			}
		}
	case DirectionDown:
		if line < e.BufferView.LineCount()-1 {
			e.window.cursor.line++
			nextLineLength := e.BufferView.LineLength(line + 1)
			if nextLineLength-1 < char {
				e.window.cursor.char = nextLineLength - 1
			}
		}
	case DirectionLeft:
		if 0 < char {
			e.window.cursor.char--
		}
	case DirectionRight:
		currentLineLength := e.BufferView.LineLength(line)
		if char < currentLineLength-1 {
			e.window.cursor.char++
		}
	default:
		panic("Invalid move direction!")
	}
}

func (e *Editor) DeleteCharBefore() {
	cursorIndex := e.getCursorPosition()
	if cursorIndex == 0 {
		return
	}
	cursorIndex--
	newLine, newChar := e.BufferView.PositionOf(cursorIndex)

	deletedChar := e.Buffer.Read(cursorIndex, cursorIndex+1)[0]
	e.Buffer.Delete(cursorIndex, cursorIndex+1)

	if deletedChar == '\n' {
		e.BufferView.Update()
	}

	e.window.cursor.line = newLine
	e.window.cursor.char = newChar
}

func (e *Editor) PutChar(c rune) {
	cursorIndex := e.getCursorPosition()
	e.Buffer.PutChar(c, cursorIndex)
	e.window.cursor.char++
}
