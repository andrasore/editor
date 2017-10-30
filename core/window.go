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
	Size() (int, int)
	Clear()
	Flush()
}

func printStatus(s Screen, statusInt, y int) {
	var statusName string

	switch statusInt {
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
		s.SetCell(x, y, c, fg, bg)
	}
}

func (w *Window) printText(s Screen, text []rune, left, top int) {
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
	}
}

func (w *Window) redraw(s Screen, b core.Buffer) {
	s.Clear()
	width, height := s.Size()

	if w.Height != height || w.Width != width {
		w.Height = height
		w.Width = width
	}

	printText(s, b.Read(0, b.Size()), 0, 0, width, height)
	printStatus(w.State, w.Height-1)
	printCursor(w.State, w.Height-1)
	s.Flush()
}

type window struct {
	width, height int
	cursor        int
}
