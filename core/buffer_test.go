package core

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := newDefaultBuffer(strings.NewReader("0123"))
	defaultEdit := edit{[]rune("0123"), 0, 4}
	assert.Equal(t, "0123", string(buf.data))
	assert.Equal(t, "", string(buf.userData))
	assert.Equal(t, []edit{defaultEdit}, buf.edits)
}

func TestBuffer_ReadMultipleEdits(t *testing.T) {
	buf := newDefaultBuffer(strings.NewReader("0123"))
	assert.Equal(t, "0123", string(buf.Read(0, 4)))
	buf.edits = append(buf.edits, edit{[]rune{'b'}, 0, 1})
	assert.Equal(t, "0123b", string(buf.Read(0, 5)))
	buf.edits = append(buf.edits, edit{[]rune{'x', 'a', 'x'}, 1, 1})
	assert.Equal(t, "0123ba", string(buf.Read(0, 6)))
}

func TestBuffer_ReadMultipleEdit_Parts(t *testing.T) {
	buf := newDefaultBuffer(strings.NewReader("0123"))
	assert.Equal(t, "0123", string(buf.Read(0, 4)))
	buf.edits = append(buf.edits, edit{[]rune("asdf"), 0, 4})
	buf.edits = append(buf.edits, edit{[]rune("xyxy"), 0, 4})
	assert.Equal(t, "3asdfxy", string(buf.Read(3, 7)))
}
