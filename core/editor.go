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

func (e *Editor) getCursorIndex() int {
	char := e.window.cursor.char
	line := e.window.cursor.line
	index, err := e.BufferView.GetCursorIndex(line, char)
	switch err.(type) {
	case IndexError:
		panic(err.Error())
	}
	return index
}

func (e *Editor) NewLine() {
	cursorIndex := e.getCursorIndex()
	e.Buffer.PutChar('\n', cursorIndex)
	e.window.cursor.char = 0
	e.window.cursor.line++
	e.BufferView.Update(0, e.Buffer.Size())
}

func (e *Editor) DeleteChar() {
	cursorIndex := e.getCursorIndex()
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
	cursorIndex := e.getCursorIndex()
	e.Buffer.PutChar(c, cursorIndex)
	e.window.cursor.char++
}
