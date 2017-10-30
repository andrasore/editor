package main

import (
	"editor/core"
	"github.com/nsf/termbox-go"
)

type termboxScreen struct{}

func (t termboxScreen) SetCell(x, y int, c rune, fg, bg int) {
}

func (t termboxScreen) Size() (int, int) {
	return termbox.Size()
}

func (t termboxScreen) Clear() {
	termbox.Clear()
}

func (t termboxScreen) Flush() {

}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	editor := core.Editor{}

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
