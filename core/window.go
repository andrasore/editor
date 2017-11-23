package core

import "fmt"

const (
	statusNormal = 0
	statusInsert = iota
)

const (
	ColorDefault = 0
	ColorBlack   = iota
	ColorWhite   = iota
)

type Screen interface {
	SetCell(x, y int, c rune, fg, bg int)
	SetCursor(x, y int)
	Size() (width int, height int)
	Clear()
	Flush()
}

func (w *window) printStatus(s Screen, status int) {
	var statusName string

	switch status {
	case statusNormal:
		statusName = "Normal"
	case statusInsert:
		statusName = "Insert"
	default:
		statusName = "???"
	}

	fg := ColorWhite
	bg := ColorBlack
	for x, c := range statusName {
		s.SetCell(x, w.height-1, c, fg, bg)
	}
}

func (w *window) printCursor(s Screen) {
	s.SetCursor(w.cursor.char, w.cursor.line)
}

func (w *window) printText(s Screen, text []rune, r rectangle) {
	fg := ColorDefault
	bg := ColorDefault

	x := r.left
	y := r.top

	for _, c := range text {
		if c == '\n' {
			x = r.left
			y++
		} else if y <= r.bottom {
			if x <= r.right {
				s.SetCell(x, y, c, fg, bg)
			}
			x++
		} else {
			break
		}
	}
}

func (w *window) debugPrintChar(s Screen, c rune) {
	r := rectangle{
		left:   w.width - 4,
		right:  w.width - 1,
		top:    w.height - 1,
		bottom: w.height - 1,
	}
	w.printText(s, []rune(fmt.Sprintf("%d", c)), r)
}

func (w *window) redraw(s Screen, b Buffer, state int, lastChar rune) {
	s.Clear()
	width, height := s.Size()

	if w.height != height || w.width != width {
		w.height = height
		w.width = width
	}

	textField := rectangle{
		left:   0,
		right:  w.width - 1,
		top:    0,
		bottom: w.height - 2,
	}
	w.printText(s, b.Read(0, b.Size()), textField)
	w.printStatus(s, state)
	w.printCursor(s)
	w.debugPrintChar(s, lastChar)
	s.Flush()
}

type rectangle struct {
	left, top, right, bottom int
}

func (r rectangle) contains(x, y int) bool {
	return r.left <= x && x <= r.right && r.bottom <= y && y <= r.top
}

type cursor struct {
	line, char int
}

type window struct {
	width, height int
	cursor        cursor
}
