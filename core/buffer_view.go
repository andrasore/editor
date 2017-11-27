package core

type BufferView interface {
	GetLineCount() int
	//GetLineLength(index int) //TODO
	GetLine(index int) []rune
	//GetLinePart(index, from, to int) []rune //TODO
	Update(from, to int)
	GetCursorPosition(line, char int) int
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

func (bw *defaultBufferView) GetLineCount() int {
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

func (bw *defaultBufferView) GetLine(index int) []rune {
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

func (bw *defaultBufferView) GetCursorPosition(line, char int) int {
	if bw.buffer.Size() == 0 {
		return 0
	}

	if bw.GetLineCount() <= line || len(bw.GetLine(line)) < char {
		panic("GetCursorPosition index out of bounds!")
	}

	if line == 0 {
		return char
	} else {
		return bw.lineIndices[line-1] + char + 1
	}
}
