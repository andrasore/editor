package test

import (
	"editor/core"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBuffer_Read(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123"))
	assert.Equal(t, "0123", string(buf.Read(0, 4)))
	assert.Equal(t, "", string(buf.Read(0, 0)))
	assert.Equal(t, "0", string(buf.Read(0, 1)))
	assert.Equal(t, "1", string(buf.Read(1, 2)))
	assert.Equal(t, "12", string(buf.Read(1, 3)))
}

func TestBuffer_InsertEmpty(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.Insert([]rune("abcd"), 0)
	result := buf.Read(0, 4)
	assert.Equal(t, "abcd", string(result))
}

func TestBuffer_Insert(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("abcd"), 0)
	assert.Equal(t, "abcd012345", string(buf.Read(0, 10)))
	buf.Insert([]rune("xx"), 2)
	assert.Equal(t, "abxxcd0123", string(buf.Read(0, 10)))
}

func TestBuffer_InsertToEnd(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("abcd"), 10)
	assert.Equal(t, "0123456789abcd", string(buf.Read(0, 14)))
	buf.Insert([]rune("x"), 13)
	assert.Equal(t, "0123456789abcxd", string(buf.Read(0, 15)))
}

func TestBuffer_DeleteEmpty(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.Delete(3, 5)
}

func TestBuffer_Delete(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Delete(3, 5)
	assert.Equal(t, "01256789", string(buf.Read(0, 8)))
}

func TestBuffer_DeleteMultiple(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("abcd"), 4) // 0123abcd456789
	buf.Delete(2, 6)
	assert.Equal(t, "01cd456789", string(buf.Read(0, 10)))
	buf.Delete(1, 5)
	assert.Equal(t, "056789", string(buf.Read(0, 10)))
}

func TestBuffer_DeleteAcrossInserts(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("123"))
	buf.Insert([]rune("456"), 3)
	buf.Insert([]rune("789"), 6)
	buf.Delete(2, 7)
	assert.Equal(t, "1289", string(buf.Read(0, 4)))
}

func TestBuffer_DeleteWholeInsert(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("123"))
	buf.Insert([]rune("456"), 3)
	buf.Delete(3, 7)
	assert.Equal(t, "123", string(buf.Read(0, 4)))
}

func TestBuffer_Size(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("abcd"), 0)
	assert.Equal(t, 14, buf.Size())
}

func TestBuffer_InsertOne(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.Insert([]rune{'a'}, 0)
	buf.Insert([]rune{'b'}, 1)
	buf.Insert([]rune{'c'}, 2)
	assert.Equal(t, "abc", string(buf.Read(0, 3)))
	assert.Equal(t, 3, buf.Size())
}

func TestBuffer_PutChar(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.PutChar('a', 0)
	assert.Equal(t, "a", string(buf.Read(0, 1)))
	buf.PutChar('b', 1)
	buf.PutChar('c', 2)
	assert.Equal(t, "abc", string(buf.Read(0, 3)))
	assert.Equal(t, 3, buf.Size())
}
