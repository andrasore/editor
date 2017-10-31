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
	switch c {
	case KeyEsc:
		if e.state == statusInsert {
			e.state = statusNormal
		}
	case 'i':
		if e.state == statusNormal {
			e.state = statusInsert
		}
	case KeyEnter:
		if e.state == statusInsert {
			e.Buffer.Insert([]rune{'\n'}, e.window.cursor)
			e.window.cursor++
		}
	case KeyBackspace:
		if e.state == statusInsert {
			e.Buffer.Delete(e.window.cursor, e.window.cursor+1)
			e.window.cursor--
		}
	default:
		if e.state == statusInsert {
			e.Buffer.PutChar(c, e.window.cursor)
			e.window.cursor++
		}
	}
	e.window.redraw(e.Screen, e.Buffer, e.state)
}
