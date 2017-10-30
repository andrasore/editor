package main

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
	assert.Equal(t, "12", string(buf.Read(1, 2)))
	assert.Equal(t, "123", string(buf.Read(1, 3)))
}

func TestBuffer_InsertEmpty(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.Insert([]rune("asdf"), 0)
	result := buf.Read(0, 4)
	assert.Equal(t, "asdf", string(result))
}

func TestBuffer_Insert(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("asdf"), 0)
	assert.Equal(t, "asdf012345", string(buf.Read(0, 10)))
	buf.Insert([]rune("xx"), 2)
	assert.Equal(t, "asxxdf0123", string(buf.Read(0, 10)))
}

func TestBuffer_Delete(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Delete(3, 2)
	assert.Equal(t, "01256789", string(buf.Read(0, 8)))
}

func TestBuffer_DeleteMultiple(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("asdf"), 4)
	buf.Delete(2, 4)
	assert.Equal(t, "01df456789", string(buf.Read(0, 10)))
	buf.Delete(1, 4)
	assert.Equal(t, "056789", string(buf.Read(0, 10)))
}

func TestBuffer_Size(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("asdf"), 0)
	assert.Equal(t, 14, buf.Size())
}

func TestBuffer_InsertOne(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.Insert([]rune{'a'}, 0)
	buf.Insert([]rune{'b'}, 1)
	buf.Insert([]rune{'c'}, 2)
	assert.Equal(t, "abc", string(buf.Read(0, 3)))
}

func TestBuffer_PutChar(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.PutChar('a', 0)
	assert.Equal(t, "a", string(buf.Read(0, 1)))
	buf.PutChar('b', 1)
	buf.PutChar('c', 2)
	assert.Equal(t, "abc", string(buf.Read(0, 3)))
}
