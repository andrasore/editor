package core

import (
	"bufio"
	"fmt"
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

func (b *defaultBuffer) Read(from, to int) (text []rune) {
	if to < from || from < 0 {
		panic(fmt.Sprintf("Trying to read from %v to %v", from, to))
	}
	if from == to {
		return
	}
	charsToRead := to - from
	editBegin := 0
	for _, e := range b.edits {
		editEnd := editBegin + len(e)
		if hasIntersection(from, to, editBegin, editEnd) {
			readFrom, readTo := intersect(from, to, editBegin, editEnd)
			readSlice := e[readFrom-editBegin : readTo-editBegin]
			text = append(text, readSlice...)
			charsToRead -= readTo - readFrom
		}
		if charsToRead == 0 {
			break
		}
		editBegin += len(e)
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

func insertIntoEdits(edits [][]rune, insertedEdit []rune, from int) [][]rune {
	var newEdits [][]rune

	editBegin := 0
	for _, e := range edits {
		if shouldSplitEdit(from, editBegin, len(e)) {
			offset := from - editBegin
			newEdits = append(newEdits, e[0:offset], insertedEdit, e[offset:])
		} else if from == editBegin {
			newEdits = append(newEdits, insertedEdit, e)
		} else {
			newEdits = append(newEdits, e)
		}
		editBegin += len(e)
	}
	if from == editBegin {
		newEdits = append(newEdits, insertedEdit)
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
		panic(fmt.Sprintf("Trying to delete from %v to %v", from, to))
	}
	if from == to {
		return
	}
	var remainingEdits [][]rune
	isDeleting := false
	editBegin := 0
	for _, e := range b.edits {
		fromInIndex := containsIndex(from, editBegin, len(e))
		toInIndex := containsIndex(to-1, editBegin, len(e))
		if fromInIndex || toInIndex {
			if fromInIndex {
				splitEdit := e[:from-editBegin]
				remainingEdits = append(remainingEdits, splitEdit)
				isDeleting = true
			}
			if toInIndex {
				splitEdit := e[to-editBegin:]
				remainingEdits = append(remainingEdits, splitEdit)
				isDeleting = false
			}
		} else if !isDeleting {
			remainingEdits = append(remainingEdits, e)
		}
		editBegin += len(e)
	}
	b.edits = remainingEdits
}

func containsIndex(index, editBegin, editLength int) bool {
	return editBegin <= index && index < editBegin+editLength
}

func shouldSplitEdit(insert, editBegin, editLength int) bool {
	return editBegin < insert && insert < editBegin+editLength-1 //last index
}

func hasIntersection(begin1, end1, begin2, end2 int) bool {
	return begin1 < end2 && begin2 < end1
}

func intersect(begin1, end1, begin2, end2 int) (int, int) {
	return max(begin1, begin2), min(end1, end2)
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
