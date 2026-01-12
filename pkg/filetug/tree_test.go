package filetug

import (
	"testing"

	"github.com/rivo/tview"
)

func TestNewTree(t *testing.T) {
	nav := NewNavigator(tview.NewApplication())
	tree := NewTree(nav)
	if tree == nil {
		t.Fatal("tree is nil")
	}
	if tree.GetRoot() == nil {
		t.Fatal("root is nil")
	}
}
