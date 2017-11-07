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
	char := e.window.cursor.char
	line := e.window.cursor.line
	cursorIndex := e.BufferView.GetCursorIndex(char, line)
	switch e.state {
	case statusInsert:
		switch c {
		case KeyEsc:
			e.state = statusNormal
		case KeyEnter:
			e.Buffer.PutChar('\n', cursorIndex)
			e.window.cursor.char = 0
			e.window.cursor.line++
		case KeyBackspace:
			e.Buffer.Delete(cursorIndex, 1)
			e.window.cursor.char-- //TODO
		default:
			e.Buffer.PutChar(c, cursorIndex)
			e.window.cursor.char++ //TODO
		}
	case statusNormal:
		switch c {
		case 'i':
			e.state = statusInsert
		}
	}

	e.window.redraw(e.Screen, e.Buffer, e.state, c)
}
