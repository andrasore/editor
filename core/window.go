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
	Size() (int, int)
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
	x := w.cursor % w.width
	y := w.cursor / w.width
	s.SetCursor(x, y)
}

func (w *window) printText(s Screen, text []rune, left, top int) {
	fg := ColorDefault
	bg := ColorDefault

	x := left
	y := top

	for _, c := range text {
		if c == '\n' {
			x = left
			y++
		} else {
			s.SetCell(x, y, c, fg, bg)
			x++
		}
	} //TODO: wrap text at width
}

func (w *window) debugPrintChar(s Screen, c rune) {
	w.printText(s, []rune(fmt.Sprintf("%d", c)), w.width-4, w.height-1)
}

func (w *window) redraw(s Screen, b Buffer, state int, lastChar rune) {
	s.Clear()
	width, height := s.Size()

	if w.height != height || w.width != width {
		w.height = height
		w.width = width
	}

	w.printText(s, b.Read(0, b.Size()), 0, 0)
	w.printStatus(s, state)
	w.printCursor(s)
	w.debugPrintChar(s, lastChar)
	s.Flush()
}

type window struct {
	width, height int
	cursor        int
}
