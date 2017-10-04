package main

import (
	"editor/core"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBuffer(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		buf := core.NewBuffer(strings.NewReader("0123"))
		assert.Equal(t, "0123", string(buf.Read(0, 4)))
		assert.Equal(t, "12", string(buf.Read(1, 2)))
		assert.Equal(t, "", string(buf.Read(0, 0)))
	})

	t.Run("Insert_Then_Read", func(t *testing.T) {
		buf := core.NewEmptyBuffer()
		buf.Insert([]rune("asdf"), 0, 0)
		result := buf.Read(0, 4)
		assert.Equal(t, "asdf", string(result))
	})

	t.Run("Insert_ToPosition", func(t *testing.T) {
		buf := core.NewBuffer(strings.NewReader("0123456789"))
		buf.Insert([]rune("asdf"), 0, 5)
		assert.Equal(t, "asdf456789", string(buf.Read(0, 10)))
	})
}
