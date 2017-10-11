package main

import (
	"editor/core"

	"github.com/nsf/termbox-go"
)

func printStatus(statusInt int) {
	//x, y := 0, 0
	//for _, c := range status.name {
	//termbox.SetCell(x, y, c, status.fg, status.bg)
	//x++
	//}
	//termbox.Flush()
}

func printText(text string, x, y int) {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	for _, c := range text {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	state := core.StatusNormal
	printStatus(state)

	//var buffer []rune = make([]rune, 1024)

	running := true
	for running {
		ev := termbox.PollEvent()
		//x, y := termbox.Size()

		switch ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyEsc:
				state = core.StatusNormal
			case ev.Ch == 'i':
				state = core.StatusInsert
			case ev.Key == termbox.KeyCtrlC:
				running = false
			}
		case termbox.EventResize:
			termbox.Flush()
		}
		printStatus(state)
	}
}
