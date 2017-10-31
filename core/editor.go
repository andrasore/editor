package core

const (
	KeyEsc       rune = 0x1B
	KeyEnter     rune = 0x0D
	KeyBackspace rune = 0x08
)

type Editor struct {
	Screen Screen
	Buffer Buffer
	window window
	state  int
}

func (e *Editor) SendChar(c rune) {
	switch e.state {
	case statusInsert:
		switch c {
		case KeyEsc:
			e.state = statusNormal
		case KeyEnter:
			e.Buffer.Insert([]rune{'\n'}, e.window.cursor)
			e.window.cursor++
		case KeyBackspace:
			e.Buffer.Delete(e.window.cursor, 1)
			e.window.cursor--
		default:
			e.Buffer.PutChar(c, e.window.cursor)
			e.window.cursor++
		}
	case statusNormal:
		switch c {
		case 'i':
			e.state = statusInsert
		}
	}

	e.window.redraw(e.Screen, e.Buffer, e.state, c)
}
