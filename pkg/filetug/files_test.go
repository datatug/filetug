package filetug

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestFiles(t *testing.T) {
	app := tview.NewApplication()
	nav := NewNavigator(app, OnMoveFocusUp(func(source tview.Primitive) {}))
	f := nav.files

	t.Run("FocusBlur", func(t *testing.T) {
		nav.filesFocusFunc()
		nav.filesBlurFunc()
	})

	t.Run("SelectionChanged", func(t *testing.T) {
		nav.filesSelectionChangedFunc(0, 0)

		f.SetCell(1, 0, tview.NewTableCell(" file.txt"))
		nav.filesSelectionChangedFunc(1, 0)

		// Test with no space prefix (should not happen in real app but for coverage)
		f.SetCell(2, 0, tview.NewTableCell("file.txt"))
		defer func() { _ = recover() }()
		nav.filesSelectionChangedFunc(2, 0)
	})

	t.Run("InputCapture_Space", func(t *testing.T) {
		f.SetCell(1, 0, tview.NewTableCell(" name"))
		f.Select(1, 0)
		event := tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone)
		f.GetInputCapture()(event)
		assert.Contains(t, f.GetCell(1, 0).Text, "✓")

		f.GetInputCapture()(event)
		assert.Contains(t, f.GetCell(1, 0).Text, " ")
	})

	t.Run("InputCapture_Keys", func(t *testing.T) {
		f.GetInputCapture()(tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone))
		f.GetInputCapture()(tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone))
		f.GetInputCapture()(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))

		// Test KeyUp at row 0
		f.Select(0, 0)
		f.GetInputCapture()(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))

		f.GetInputCapture()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
	})
}
