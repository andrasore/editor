package core

const (
	KeyEsc       rune = 0x1B
	KeyEnter     rune = 0x0D
	KeyBackspace rune = 0x08
)

type Editor struct {
	Screen     Screen
	Buffer     Buffer
	BufferView BufferView
	window     window
	state      int
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
			e.DeleteChar()
		default:
			e.PutChar(c)
		}
	case statusNormal:
		switch c {
		case 'i':
			e.state = statusInsert
		}
	}

	e.window.redraw(e.Screen, e.Buffer, e.state, c)
}

func (e *Editor) getCursorPosition() int {
	char := e.window.cursor.char
	line := e.window.cursor.line
	return e.BufferView.GetCursorPosition(line, char)
}

func (e *Editor) NewLine() {
	cursorPosition := e.getCursorPosition()
	e.Buffer.PutChar('\n', cursorPosition)
	e.window.cursor.char = 0
	e.window.cursor.line++
	e.BufferView.Update(0, e.Buffer.Size()) //TODO
}

func (e *Editor) DeleteChar() {
	cursorIndex := e.getCursorPosition()
	if cursorIndex == 0 {
		return
	}

	e.Buffer.Delete(cursorIndex, 1)

	if e.window.cursor.char != 0 {
		e.window.cursor.char--
	} else {
		e.window.cursor.line--
		line := e.BufferView.GetLine(e.window.cursor.line)
		e.window.cursor.char = len(line) - 1
	}
}

func (e *Editor) PutChar(c rune) {
	cursorIndex := e.getCursorPosition()
	e.Buffer.PutChar(c, cursorIndex)
	e.window.cursor.char++
}
