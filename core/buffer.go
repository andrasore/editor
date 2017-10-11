package core

import (
	"bufio"
	"io"
	"strings"
)

type Buffer interface {
	Insert(text []rune, from, to int)
	Read(from, to int) []rune
	Size() int
}

type defaultBuffer struct {
	data     []rune
	userData []rune
	edits    []edit
}

func newDefaultBuffer(reader io.Reader) *defaultBuffer {
	bufioReader := bufio.NewReader(reader)
	content := make([]rune, 0)

	r, _, err := bufioReader.ReadRune()

	for err == nil {
		content = append(content, r)
		r, _, err = bufioReader.ReadRune()
	}

	if err != io.EOF {
		panic("Todo: implement error handling")
	}

	initialEdit := edit{content, 0, len(content)}

	return &defaultBuffer{
		content,
		make([]rune, 0),
		[]edit{initialEdit},
	}
}

func NewBuffer(reader io.Reader) Buffer {
	return newDefaultBuffer(reader)
}

func NewEmptyBuffer() Buffer {
	return NewBuffer(strings.NewReader(""))
}

func isIndexInEdit(index, indexInText, editLength int) bool {
	return index < indexInText+editLength
}

func (b *defaultBuffer) Read(from, length int) (text []rune) {
	indexInText := 0
	charsToRead := length
	for _, e := range b.edits {
		if isIndexInEdit(from, indexInText, e.length) {
			readOffset := max(from, indexInText) - indexInText
			readBegin := e.dataIndex + readOffset
			currentReadSize := min(e.length-readOffset, charsToRead)
			text = append(text, e.data[readBegin:readBegin+currentReadSize]...)
			charsToRead -= currentReadSize
		}
		if charsToRead == 0 {
			break
		}
		indexInText += e.length
	}
	return
}

func (b *defaultBuffer) Insert(text []rune, from, to int) {
	editBegin := len(b.userData)
	b.userData = append(b.userData, text...)

	newEdit := edit{b.userData, editBegin, len(text)}

	b.edits = getNewEdits(*b, newEdit, from, to)
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

func (b defaultBuffer) iterateEdits(iterFunc func(int, edit) bool) {
	previousLength := 0
	for _, e := range b.edits {
		if iterFunc(previousLength, e) {
			break
		}
		previousLength += e.length
	}
}

func (b defaultBuffer) Size() int {
	return 0
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

type edit struct {
	data              []rune
	dataIndex, length int
}
