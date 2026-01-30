package filetug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowHelpModal(t *testing.T) {
	nav, _, _ := newNavigatorForTest(t)
	showHelpModal(nav)
	assert.NotNil(t, nav)
}
