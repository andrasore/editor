package main

import (
	"editor/core"
	"github.com/nsf/termbox-go"
)

func getTermboxColor(coreColor int) termbox.Attribute {
	switch coreColor {
	case core.ColorBackground:
		return termbox.ColorWhite
	case core.ColorForeground:
		return termbox.ColorBlack
	default:
		return termbox.ColorDefault
	}
}

type termboxScreen struct{}

func (t termboxScreen) SetCell(x, y int, c rune, fg, bg int) {
	termbox.SetCell(x, y, c, getTermboxColor(fg), getTermboxColor(bg))
}

func (t termboxScreen) Size() (int, int) {
	return termbox.Size()
}

func (t termboxScreen) Clear() {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	termbox.Clear(fg, bg)
}

func (t termboxScreen) Flush() {
	termbox.Flush()
}

func (t termboxScreen) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	editor := core.Editor{Screen: termboxScreen{}, Buffer: core.NewEmptyBuffer()}

	editor.SendChar(core.KeyEsc) //TODO: redraw nicely

	running := true
	for running {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			editor.SendChar(ev.Ch)
		}
		if ev.Key == termbox.KeyCtrlC {
			running = false
		}
	}
}
