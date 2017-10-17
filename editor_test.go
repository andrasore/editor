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

func TestBuffer_Insert(t *testing.T) {
	buf := core.NewEmptyBuffer()
	buf.Insert([]rune("asdf"), 0, 0)
	result := buf.Read(0, 4)
	assert.Equal(t, "asdf", string(result))
}

func TestBuffer_InsertOverwrite(t *testing.T) {
	buf := core.NewBuffer(strings.NewReader("0123456789"))
	buf.Insert([]rune("asdf"), 0, 5)
	assert.Equal(t, "asdf56789", string(buf.Read(0, 10)))
	buf.Insert([]rune("xx"), 2, 0)
	assert.Equal(t, "asxxdf5678", string(buf.Read(0, 10)))
}
