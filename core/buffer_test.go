package core

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := newDefaultBuffer(strings.NewReader("0123"))
	defaultEdit := []rune("0123")
	assert.Equal(t, "", string(buf.userData))
	assert.Equal(t, [][]rune{defaultEdit}, buf.edits)
}

func TestBuffer_ReadMultipleEdits(t *testing.T) {
	buf := newDefaultBuffer(strings.NewReader("0123"))
	assert.Equal(t, "0123", string(buf.Read(0, 4)))
	buf.edits = append(buf.edits, []rune{'b'})
	assert.Equal(t, "0123b", string(buf.Read(0, 5)))
	buf.edits = append(buf.edits, []rune{'a'})
	assert.Equal(t, "0123ba", string(buf.Read(0, 6)))
}

func TestBuffer_ReadMultipleEdit_Parts(t *testing.T) {
	buf := newDefaultBuffer(strings.NewReader("0123"))
	assert.Equal(t, "0123", string(buf.Read(0, 4)))
	buf.edits = append(buf.edits, []rune("asdf"))
	buf.edits = append(buf.edits, []rune("xyxy"))
	assert.Equal(t, "3asdfxy", string(buf.Read(3, 7)))
}
