package filetug

import (
	"context"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestTree(t *testing.T) {
	app := tview.NewApplication()
	nav := NewNavigator(app)
	tree := NewTree(nav)

	t.Run("onStoreChange", func(t *testing.T) {
		tree.onStoreChange()
	})

	t.Run("Draw", func(t *testing.T) {
		screen := tcell.NewSimulationScreen("")
		_ = screen.Init()
		tree.Draw(screen)
	})

	t.Run("changed", func(t *testing.T) {
		root := tree.GetRoot()
		tree.changed(root)
	})

	t.Run("focus_blur", func(t *testing.T) {
		tree.focus()
		tree.blur()
	})

	t.Run("inputCapture", func(t *testing.T) {
		eventLeft := tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
		tree.inputCapture(eventLeft)

		eventRight := tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
		tree.inputCapture(eventRight)

		eventUp := tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		tree.inputCapture(eventUp)
	})

	t.Run("SetSearch", func(t *testing.T) {
		tree.SetSearch("test")
	})

	t.Run("setCurrentDir", func(t *testing.T) {
		tree.setCurrentDir("/")
	})

	t.Run("setDirContext", func(t *testing.T) {
		root := tree.GetRoot()
		dc := &DirContext{Path: "/test"}
		tree.setDirContext(context.Background(), root, dc)
	})

	t.Run("setError", func(t *testing.T) {
		root := tree.GetRoot()
		tree.setError(root, context.DeadlineExceeded)
	})

	t.Run("getNodePath", func(t *testing.T) {
		root := tree.GetRoot()
		root.SetReference("/")
		child := tview.NewTreeNode("child")
		child.SetReference("/child")
		root.AddChild(child)
		getNodePath(child)
	})
}
