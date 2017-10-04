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
	content := make([]rune, 0)

	r, _, err := bufioReader.ReadRune()

	for err == nil {
		content = append(content, r)
		r, _, err = bufioReader.ReadRune()
	}

	if err != io.EOF {
		panic("Buffer done fucked up")
	}

	initialEdit := edit{content, 0, len(content)}

	return &defaultBuffer{
		content,
		make([]rune, 0),
		[]edit{initialEdit},
	}
}

func NewEmptyBuffer() Buffer {
	return NewBuffer(strings.NewReader(""))
}

func (b *defaultBuffer) Insert(text []rune, from, to int) {
	b.userData = append(b.userData, text...)

	newEdit := edit{b.userData, len(b.userData), len(text)}

	b.edits = getNewEdits(*b, newEdit, from, to)
}

func (b *defaultBuffer) Read(from, length int) (text []rune) {
	toRead := length
	b.iterateEdits(func(textIndex int, e edit) bool {
		if textIndex <= from {
			if e.length < toRead {
				text = append(text, e.data[e.from:e.from+e.length]...)
				toRead -= e.length
			} else {
				text = append(text, e.data[e.from:e.from+toRead]...)
				return true
			}
		}
		return false
	})
	return
}

func (b *defaultBuffer) iterateEdits(iterFunc func(int, edit) bool) {
	previousLength, currentLength := 0, 0
	for _, e := range b.edits {
		previousLength = currentLength
		currentLength += e.length
		if iterFunc(previousLength, e) {
			break
		}
	}
}

func getNewEdits(b defaultBuffer, pastedEdit edit, from, to int) []edit {
	var newEdits []edit

	b.iterateEdits(func(textIndex int, e edit) bool {
		if from < textIndex+e.length {
			splitEdit := edit{e.data, textIndex, from - textIndex}
			newEdits = append(newEdits, splitEdit, pastedEdit)
		} else {
			newEdits = append(newEdits, e)
		}

		if to <= textIndex+e.length {
			return true
		}
		return false
	})

	return newEdits
}

type edit struct {
	data         []rune
	from, length int
}
