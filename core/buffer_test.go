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
