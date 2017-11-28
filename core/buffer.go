package core

import (
	"bufio"
	"io"
	"strings"
)

type Buffer interface {
	Insert(text []rune, from int)
	PutChar(char rune, from int)
	Delete(from, to int)
	Read(from, to int) []rune
	Size() int
}

type defaultBuffer struct {
	userData []rune
	edits    [][]rune
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

	var edits [][]rune
	if len(content) != 0 {
		edits = [][]rune{content}
	}

	return &defaultBuffer{make([]rune, 0), edits}
}

func NewBuffer(reader io.Reader) Buffer {
	return newDefaultBuffer(reader)
}

func NewEmptyBuffer() Buffer {
	return NewBuffer(strings.NewReader(""))
}

func hasIntersection(begin1, end1, begin2, length2 int) bool {
	end2 := begin2 + length2
	return begin1 < end2 && begin2 < end1
}

func (b *defaultBuffer) Read(from, to int) (text []rune) {
	if to < from || from < 0 {
		panic("Invalid indices for read!")
	}
	charsToRead := to - from
	indexInText := 0
	for _, e := range b.edits { //TODO: implement a linked list
		if hasIntersection(from, to, indexInText, len(e)) {
			readOffset := max(from, indexInText) - indexInText
			currentReadSize := min(len(e)-readOffset, charsToRead)
			text = append(text, e[readOffset:readOffset+currentReadSize]...)
			charsToRead -= currentReadSize
		}
		if charsToRead == 0 {
			break
		}
		indexInText += len(e)
	}
	return
}

func (b *defaultBuffer) PutChar(char rune, from int) {
	b.Insert([]rune{char}, from)
} //TODO: this sucks! - remembering last insert position would be nice

func (b *defaultBuffer) Insert(text []rune, from int) {
	if len(text) == 0 {
		return
	}

	editBegin := len(b.userData)
	b.userData = append(b.userData, text...)
	newEdit := b.userData[editBegin:]

	if len(b.edits) == 0 {
		b.edits = [][]rune{newEdit}
	} else {
		b.edits = insertIntoEdits(b.edits, newEdit, from)
	}
}

func insertIntoEdits(edits [][]rune, pastedEdit []rune, from int) [][]rune {
	var newEdits [][]rune

	indexInText := 0
	for _, e := range edits {
		if containsIndex(from, indexInText, len(e)) {
			newEditOffset := from - indexInText

			firstSplitEdit := e[0:newEditOffset]
			if len(firstSplitEdit) > 0 {
				newEdits = append(newEdits, firstSplitEdit)
			}

			newEdits = append(newEdits, pastedEdit)

			secondSplitEdit := e[newEditOffset:]
			if len(secondSplitEdit) > 0 {
				newEdits = append(newEdits, secondSplitEdit)
			}
		} else {
			newEdits = append(newEdits, e)
		}

		indexInText += len(e)
	}

	return newEdits
}

func (b *defaultBuffer) Size() int {
	size := 0
	for _, e := range b.edits {
		size += len(e)
	}
	return size
}

func (b *defaultBuffer) Delete(from, to int) {
	if to < from || from < 0 {
		panic("Invalid indices for delete!")
	}

	b.edits = getRemainingEdits(b.edits, from, to)
}

func getRemainingEdits(edits [][]rune, delFrom, delTo int) [][]rune {
	var remainingEdits [][]rune
	isDeleting := false
	textIndex := 0
	for _, e := range edits {
		fromInIndex := containsIndex(delFrom, textIndex, len(e))
		toInIndex := containsIndex(delTo-1, textIndex, len(e))
		if fromInIndex || toInIndex {
			if fromInIndex {
				splitEdit := e[:delFrom-textIndex]
				remainingEdits = append(remainingEdits, splitEdit)
				isDeleting = true
			}
			if toInIndex {
				splitEdit := e[delTo-textIndex:]
				remainingEdits = append(remainingEdits, splitEdit)
				isDeleting = false
			}
		} else if !isDeleting {
			remainingEdits = append(remainingEdits, e)
		}
		textIndex += len(e)
	}
	return remainingEdits
}

func containsIndex(index, editBegin, editLength int) bool {
	return editBegin <= index && index < editBegin+editLength
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
