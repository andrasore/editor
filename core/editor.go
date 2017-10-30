package core

const (
	KeyEsc       rune = 0x1B
	KeyEnter     rune = 0x0D
	KeyBackspace rune = 0x08
)

type Editor struct {
	screen Screen
	buffer Buffer
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
			e.buffer.Insert([]rune{'\n'}, e.window.cursor)
			e.window.cursor++
		}
	case KeyBackspace:
		if e.state == statusInsert {
			e.buffer.Delete(e.window.cursor, e.window.cursor+1)
			e.window.cursor--
		}
	default:
		if window.state == statusInsert {
			e.buffer.PutChar(c, e.window.cursor)
			window.cursor++
		}
	}
	e.redraw()
}
