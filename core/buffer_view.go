package core

import "fmt"

type BufferView interface {
	LineCount() int
	LineLength(index int) int
	Line(index int) []rune
	//LinePart(index, from, to int) []rune //TODO
	//Text(fromLine, fromChar, toLine, toChar int) []rune
	Update(from, to int)
	IndexOf(line, char int) int
	PositionOf(index int) (line, char int)
}

type defaultBufferView struct {
	buffer      Buffer
	lineIndices []int
}

func NewBufferView(buffer Buffer) BufferView {
	bufferView := defaultBufferView{buffer: buffer}
	bufferView.Update(0, buffer.Size())
	return BufferView(&bufferView)
}

func (bw *defaultBufferView) LineCount() int {
	return len(bw.lineIndices) + 1
}

func (bw *defaultBufferView) Update(from, to int) {
	if from != 0 || to != bw.buffer.Size() {
		panic("TODO - can only be called with full range for now")
	}
	bw.lineIndices = nil
	for i, r := range bw.buffer.Read(from, to) {
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

func (bw *defaultBufferView) PositionOf(index int) (int, int) {
	if index == 0 {
		return 0, 0
	}

	if index < 0 || bw.buffer.Size() < index {
		panic(fmt.Sprintf("Index out of bounds for buffer: %v", index))
	}

	line := func() int {
		for i := 0; i < len(bw.lineIndices); i++ {
			if index < bw.lineIndices[i] {
				return i
				break
			}
		}
		return len(bw.lineIndices)
	}()

	char := index - bw.IndexOf(line, 0)

	return line, char
}
