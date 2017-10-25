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

	fg := termbox.ColorBlack
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

func redraw(w *core.Window, b *core.Buffer) {
	height, width := termbox.Size()

	if w.Height != height || w.Width != width {
		w.Height = height
		w.Width = width
		//TODO - resize?
	}

	printStatus(w.State, w.Height-1)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	window := core.Window{}
	buffer := core.NewEmptyBuffer()
	cursor := 0

	running := true
	for running {
		ev := termbox.PollEvent()

		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyEsc && window.State == core.StatusInsert:
				window.State = core.StatusNormal
			case ev.Ch == 'i' && window.State == core.StatusNormal:
				window.State = core.StatusInsert
			case ev.Key == termbox.KeyCtrlC:
				running = false
			default:
				buffer.Insert([]rune{ev.Ch}, cursor)
				cursor++
			}

		case termbox.EventResize:
			termbox.Flush()
		}

		redraw(&window, &buffer)
	}
}
