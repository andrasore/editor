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

func hasIntersection(begin1, length1, begin2, length2 int) bool {
	end1 := begin1 + length1
	end2 := begin2 + length2
	return begin1 < end2 && begin2 < end1
}

func (b *defaultBuffer) Read(from, length int) (text []rune) {
	indexInText := 0
	charsToRead := length
	for _, e := range b.edits {
		if hasIntersection(from, length, indexInText, e.length) {
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

func isIndexInEdit(index, editBegin, editLength int) bool {
	return editBegin <= index && index < editBegin+editLength
}

func getNewEdits(b defaultBuffer, pastedEdit edit, from, to int) []edit {
	var newEdits []edit

	indexInText := 0
	for _, e := range b.edits {
		containsFrom := isIndexInEdit(from, indexInText, e.length)
		containsTo := isIndexInEdit(to, indexInText, e.length)

		switch {
		case containsFrom && containsTo:
			firstSplitEdit := edit{e.data, indexInText, from - indexInText}
			secondSplitEdit := edit{e.data, to, (indexInText + e.length) - to}
			newEdits = append(newEdits, firstSplitEdit, pastedEdit, secondSplitEdit)
		case containsFrom:
			splitEdit := edit{e.data, indexInText, from - indexInText}
			newEdits = append(newEdits, splitEdit, pastedEdit)
		case containsTo:
			splitEdit := edit{e.data, to, (indexInText + e.length) - to}
			newEdits = append(newEdits, splitEdit, pastedEdit)
		default:
			newEdits = append(newEdits, e)
		}

		indexInText += e.length
	}

	return newEdits
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
