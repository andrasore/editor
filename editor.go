package main

import (
	"github.com/nsf/termbox-go"
)

const (
	StateNormal = iota
	StateInsert = iota
)

func getStatusString(status int) string {
	switch status {
	case StateNormal:
		return "Normal"
	case StateInsert:
		return "Insert"
	default:
		return "Invalid"
	}
}

func printStatus(status int) {
	statusString := getStatusString(status)
	x, y := 0, 0
	fg := termbox.ColorBlack
	bg := termbox.ColorDefault
	for _, c := range statusString {
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

	state := StateNormal
	printStatus(state)

	running := true
	for running {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyEsc:
				state = StateNormal
			case ev.Ch == 'i':
				state = StateInsert
			case ev.Key == termbox.KeyCtrlC:
				running = false
			}
		}
		printStatus(state)
	}
}
