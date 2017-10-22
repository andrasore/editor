package core

import (
	"bufio"
	"io"
	"strings"
)

type Buffer interface {
	Insert(text []rune, from int)
	Delete(from, length int)
	Read(from, length int) []rune
	Size() int
}

type defaultBuffer struct {
	data     []rune
	userData []rune
	edits    []edit //TODO - make this a linked list instead
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

	var edits []edit
	if len(content) != 0 {
		edits = []edit{initialEdit}
	}

	return &defaultBuffer{
		content,
		make([]rune, 0),
		edits,
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

func (b *defaultBuffer) Insert(text []rune, from int) {
	editBegin := len(b.userData)
	b.userData = append(b.userData, text...)
	newEdit := edit{b.userData, editBegin, len(text)}

	if len(b.edits) == 0 {
		b.edits = []edit{newEdit}
	} else {
		b.edits = insertIntoEdits(b.edits, newEdit, from)
	}
}

func isIndexInEdit(index, editBegin, editLength int) bool {
	return editBegin <= index && index < editBegin+editLength
}

func insertIntoEdits(edits []edit, pastedEdit edit, from int) []edit {
	var newEdits []edit

	indexInText := 0
	for _, e := range edits {
		if isIndexInEdit(from, indexInText, e.length) {
			newEditOffset := from - indexInText

			firstSplitEdit := edit{
				e.data,
				indexInText,
				newEditOffset,
			}

			secondSplitEdit := edit{
				e.data,
				indexInText + newEditOffset,
				e.length - newEditOffset,
			}

			newEdits = append(newEdits, firstSplitEdit, pastedEdit, secondSplitEdit)
		} else {
			newEdits = append(newEdits, e)
		}

		indexInText += e.length
	}

	return newEdits
}

func (b defaultBuffer) Size() int {
	return 0
}

func (b defaultBuffer) Delete(from, length int) {

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
