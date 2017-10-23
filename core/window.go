package core

const (
	StatusNormal = 0
	StatusInsert = iota
)

type Cursor struct {
	X, Y int
}

type Window struct {
	Width, Height int
	Cursor        Cursor
	StatusBar     []rune
	Status        int
}
