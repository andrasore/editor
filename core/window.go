package core

const (
	statusNormal = 0
	statusInsert = iota
)

const (
	ColorForeground = 0
	ColorBackground = iota
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

	fg := ColorForeground
	bg := ColorBackground
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
	fg := ColorForeground
	bg := ColorBackground

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

func (w *window) redraw(s Screen, b Buffer, state int) {
	s.Clear()
	width, height := s.Size()

	if w.height != height || w.width != width {
		w.height = height
		w.width = width
	}

	w.printText(s, b.Read(0, b.Size()), 0, 0)
	w.printStatus(s, state)
	w.printCursor(s)
	s.Flush()
}

type window struct {
	width, height int
	cursor        int
}
