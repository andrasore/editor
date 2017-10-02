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

type edit struct {
	from, to int
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

	return &defaultBuffer{
		content,
		make([]rune, 1024),
		make([]edit, 0),
	}
}

func NewEmptyBuffer() Buffer {
	return NewBuffer(strings.NewReader(""))
}

func (b *defaultBuffer) Insert(text []rune, from, to int) {
}

func (b *defaultBuffer) Read(from, to int) []rune {
	return make([]rune, to-from)
}
