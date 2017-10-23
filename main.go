package main

import (
	"editor/core"
	"github.com/nsf/termbox-go"
)

func printStatus(statusInt, y int) {
	var statusName string

	switch statusInt {
	case core.StatusNormal:
		statusName = "Normal"
	case core.StatusInsert:
		statusName = "Insert"
	default:
		statusName = "???"
	}

	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	for x, c := range statusName {
		termbox.SetCell(x, y, c, fg, bg)
	}
	termbox.Flush()
}

func printText(text []rune, x, y int) {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	for _, c := range text {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
	termbox.Flush()
}

func redraw(w *core.Window) {
	height, width := termbox.Size()

	if w.Height != height || w.Width != width {
		w.Height = height
		w.Width = width
		//TODO - resize?
	}

}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	state := StatusNormal
	printStatus(state)

	window := core.Window{}
	var buffer []rune

	running := true
	for running {
		ev := termbox.PollEvent()

		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyEsc:
				state = StatusNormal
			case ev.Ch == 'i':
				state = StatusInsert
			case ev.Key == termbox.KeyCtrlC:
				running = false
			}

		case termbox.EventResize:
			termbox.Flush()
		}

		redraw(window)
	}
}
