package sneatv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTabs(t *testing.T) {
	tabs := NewTabs(nil, UnderlineTabsStyle)
	assert.NotNil(t, tabs)
}
