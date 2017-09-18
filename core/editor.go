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

type Buffer struct {
	data []rune
}

func NewBuffer() *Buffer {
	return &Buffer{make([]rune, 0)}
}

func (b *Buffer) Insert(text []rune, from, to int) {
	copy(b.data[from:to], text)
}

func (b *Buffer) Read(from, to int) []rune {
	return b.data[from:to]
}
