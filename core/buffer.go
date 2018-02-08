package core

import (
	"bufio"
	"container/list"
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

type editListBuffer struct {
	data       []rune
	edits      list.List
	lastEdit   *list.Element
	lastInsert int
}

func NewBuffer(reader io.Reader) Buffer {
	newBuffer := editListBuffer{
		data:       make([]rune, 0),
		lastInsert: 0,
		lastEdit:   nil,
	}

	bufioReader := bufio.NewReader(reader)
	content := make([]rune, 0)

	r, _, err := bufioReader.ReadRune()

	for err == nil {
		content = append(content, r)
		r, _, err = bufioReader.ReadRune()
	}

	if err != io.EOF {
		return Buffer(&newBuffer)
	}

	if len(content) != 0 {
		newBuffer.Insert(content, 0)
	}

	return Buffer(&newBuffer)
}

func NewEmptyBuffer() Buffer {
	return NewBuffer(strings.NewReader(""))
}

func (b *editListBuffer) Read(from, to int) (text []rune) {
	if to < from || from < 0 {
		panic(fmt.Sprintf("Trying to read from %v to %v", from, to))
	}
	if from == to {
		return
	}
	charsToRead := to - from
	editBegin := 0
	for e := b.edits.Front(); e != nil; e = e.Next() {
		edit := e.Value.([]rune)
		checkZeroLenEdit(edit)
		editEnd := editBegin + len(edit)
		if hasIntersection(from, to, editBegin, editEnd) {
			readFrom, readTo := intersect(from, to, editBegin, editEnd)
			readSlice := edit[readFrom-editBegin : readTo-editBegin]
			text = append(text, readSlice...)
			charsToRead -= readTo - readFrom
		}
		if charsToRead == 0 {
			break
		}
		editBegin += len(edit)
	}
	return
}

func (b *editListBuffer) PutChar(char rune, from int) {
	b.Insert([]rune{char}, from)
}

func (b *editListBuffer) lastEditEnd() int {
	return b.lastInsert + len(b.lastEdit.Value.([]rune))
}

func (b *editListBuffer) Insert(text []rune, from int) {
	if len(text) == 0 {
		return
	}
	editBegin := len(b.data)
	b.data = append(b.data, text...)

	newEdit := b.data[editBegin:]

	if b.lastEdit != nil && b.lastEditEnd() == from {
		b.lastEdit.Value = append(b.lastEdit.Value.([]rune), newEdit...)
	} else {
		b.insertIntoEdits(newEdit, from)
	}

	b.lastInsert = from
}

func (b *editListBuffer) insertIntoEdits(insertedEdit []rune, from int) {
	editBegin := 0
	for e := b.edits.Front(); e != nil; e = e.Next() {
		edit := e.Value.([]rune)
		checkZeroLenEdit(edit)
		editEnd := editBegin + len(edit)
		if shouldSplitEdit(from, editBegin, editEnd) {
			offset := from - editBegin
			b.edits.InsertBefore(edit[0:offset], e)
			b.lastEdit = b.edits.InsertBefore(insertedEdit, e)
			e.Value = edit[offset:]
		} else if from == editBegin {
			b.lastEdit = b.edits.InsertBefore(insertedEdit, e)
		}
		editBegin += len(edit)
	}
	if from == editBegin { //buffer end actually
		b.lastEdit = b.edits.PushBack(insertedEdit)
	}
}

func (b *editListBuffer) Size() int {
	size := 0
	for e := b.edits.Front(); e != nil; e = e.Next() {
		size += len(e.Value.([]rune))
	}
	return size
}

func (b *editListBuffer) Delete(from, to int) {
	if to < from || from < 0 {
		panic(fmt.Sprintf("Trying to delete from %v to %v", from, to))
	}
	if from == to {
		return
	}
	editBegin := 0
	for e := b.edits.Front(); e != nil; e = e.Next() {
		edit := e.Value.([]rune)
		checkZeroLenEdit(edit)
		editEnd := editBegin + len(edit)
		if hasIntersection(from, to, editBegin, editEnd) {
			currentFrom, currentTo := intersect(from, to, editBegin, editEnd)
			fromOffset, toOffset := currentFrom-editBegin, currentTo-editBegin
			if 0 < fromOffset {
				newEdits.InsertBefore(edit[0:fromOffset], e)
			}
			if toOffset < len(edit) {
				newEdits.InsertBefore(edit[toOffset:], e)
			} //TODO
		}
		editBegin += len(edit)
	}
}

func shouldSplitEdit(splitAt, editBegin, editEnd int) bool {
	return editBegin < splitAt && splitAt < editEnd
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

func checkZeroLenEdit(e []rune) {
	if len(e) == 0 {
		panic("zero length edit!")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
