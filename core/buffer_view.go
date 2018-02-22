package core

import "fmt"

// BufferView provides methods to access a buffer's lines. It should be updated
// when buffer content is changed.
// Newlines are by definition the last characters of a line.
type BufferView interface {
	LineCount() int
	LineLength(index int) int
	Line(index int) []rune
	Update()
	IndexOf(line, char int) int
	PositionOf(index int) (line, char int)
}

type defaultBufferView struct {
	buffer      Buffer
	lineIndices []int
}

func NewBufferView(buffer Buffer) BufferView {
	bufferView := defaultBufferView{buffer: buffer}
	bufferView.Update()
	return BufferView(&bufferView)
}

func (bw *defaultBufferView) LineCount() int {
	return len(bw.lineIndices) + 1
}

func (bw *defaultBufferView) Update() {
	bw.lineIndices = nil
	for i, r := range bw.buffer.Read(0, bw.buffer.Size()) {
		if r == '\n' {
			bw.lineIndices = append(bw.lineIndices, i)
		}
	}
}

func (bw *defaultBufferView) Line(index int) []rune {
	lastLineIndex := len(bw.lineIndices)
	var from, to int

	switch {
	case index == 0 && lastLineIndex == 0:
		from = 0
		to = bw.buffer.Size()
	case index == 0:
		from = 0
		to = bw.lineIndices[0]
	case index < lastLineIndex:
		from = bw.lineIndices[index-1] + 1
		to = bw.lineIndices[index]
	case index == lastLineIndex:
		from = bw.lineIndices[index-1] + 1
		to = bw.buffer.Size()
	default:
		panic("GetLine index out of bounds!")
	}

	return bw.buffer.Read(from, to)
}

func (bw *defaultBufferView) LineLength(index int) int {
	return len(bw.Line(index))
}

// IndexOf returns the buffer index of the selected character.
func (bw *defaultBufferView) IndexOf(line, char int) int {
	if bw.buffer.Size() == 0 {
		return 0
	}

	if bw.LineCount() <= line || bw.LineLength(line) < char {
		panic(fmt.Sprintf("Out of bounds cursor: line %v, char %v", line, char))
	}

	if line == 0 {
		return char
	} else {
		return bw.lineIndices[line-1] + char + 1
	}
}

// PositionOf returns the line and char number of a given buffer index.
func (bw *defaultBufferView) PositionOf(index int) (line, char int) {
	if index == 0 {
		return 0, 0
	}

	if index < 0 || bw.buffer.Size() < index {
		panic(fmt.Sprintf("Index out of bounds for buffer: %v", index))
	}

	line = func() int {
		for i, line := range bw.lineIndices {
			if index <= line {
				return i
			}
		}
		return len(bw.lineIndices)
	}()

	char = index - bw.IndexOf(line, 0)
	return
}
