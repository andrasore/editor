package core

const (
	StatusNormal = iota
	StatusInsert = iota
)

type Window struct {
	width, height    int
	cursorX, cursorY int
	buffer           *Buffer
}
