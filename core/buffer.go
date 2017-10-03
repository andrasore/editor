package core

import (
	"bufio"
	"io"
	"strings"
)

type Buffer interface {
	Insert(text []rune, from, to int)
	Read(from, to int) []rune
}

type defaultBuffer struct {
	data     []rune
	userData []rune
	edits    []edit
}

func NewBuffer(reader io.Reader) Buffer {
	bufioReader := bufio.NewReader(reader)
	content := make([]rune, 1024)

	r, _, err := bufioReader.ReadRune()

	for err == nil {
		content = append(content, r)
		r, _, err = bufioReader.ReadRune()
	}

	if err != io.EOF {
		panic("Buffer done fucked up")
	}

	initialEdit := edit{0, len(content), content}

	return &defaultBuffer{
		content,
		make([]rune, 1024),
		[]edit{initialEdit},
	}
}

func NewEmptyBuffer() Buffer {
	return NewBuffer(strings.NewReader(""))
}

func (b *defaultBuffer) Insert(text []rune, from, to int) {
	b.userData = append(b.userData, text...)
	newEdit := edit{len(b.userData) - 1, len(text), b.userData}
	applyNewEdit(b.edits, newEdit, from, to)
}

func (b *defaultBuffer) Read(from, to int) []rune {
	return make([]rune, to-from)
}

func applyNewEdit(edits []edit, e edit, from int, to int) {
	currentLength := 0
	var affectedEditIndices []int

	for i, e := range edits {
		currentLength += e.length
		if from <= currentLength {
			affectedEditIndices = append(affectedEditIndices, i)
		}
	}
}

type edit struct {
	from, length int
	data         []rune
}
