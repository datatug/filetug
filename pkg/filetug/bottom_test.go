package filetug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_bottom_getCtrlMenuItems(t *testing.T) {
	b := &bottom{}
	menuItems := b.getCtrlMenuItems()
	assert.Len(t, menuItems, 4)
}
