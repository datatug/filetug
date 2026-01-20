package filetug

import (
	"os"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNewDirSummary(t *testing.T) {
	app := tview.NewApplication()
	nav := NewNavigator(app)
	ds := newDirSummary(nav)
	assert.NotNil(t, ds)
	assert.NotNil(t, ds.extTable)
}

func TestDirSummary_SetDir(t *testing.T) {
	app := tview.NewApplication()
	nav := NewNavigator(app)
	ds := newDirSummary(nav)

	entries := []os.DirEntry{
		mockDirEntry{name: "image1.png", isDir: false},
		mockDirEntry{name: "image2.png", isDir: false},
		mockDirEntry{name: "script.go", isDir: false},
		mockDirEntry{name: "unknown.foo", isDir: false},
		mockDirEntry{name: "subdir", isDir: true},
	}

	dir := &DirContext{
		Path:     "/test",
		children: entries,
	}

	ds.SetDir(dir)

	// .png -> Image, .go -> Code, .foo -> Other
	assert.Len(t, ds.extGroups, 3)

	var imageGroup *extensionsGroup
	for _, g := range ds.extGroups {
		if g.id == "Image" {
			imageGroup = g
			break
		}
	}
	if imageGroup == nil {
		t.Fatal("expected imageGroup to be not nil")
	}
	assert.Equal(t, "Images", imageGroup.title)
	assert.Len(t, imageGroup.extStats, 1) // .png
}

func TestGetSizeCell(t *testing.T) {
	testCases := []struct {
		size int64
	}{
		{1024 * 1024 * 1024 * 1024 * 2},
		{1024 * 1024 * 1024 * 2},
		{1024 * 1024 * 2},
		{1024 * 2},
		{512},
		{0},
	}

	for _, tc := range testCases {
		cell := getSizeCell(tc.size, 0)
		assert.NotEmpty(t, cell.Text)
	}
}

func TestDirSummary_Extra(t *testing.T) {
	app := tview.NewApplication()
	nav := NewNavigator(app)
	if nav == nil {
		t.Fatal("expected navigator to be not nil")
	}
	nav.files = newFiles(nav) // Ensure nav.files is initialized to avoid panic
	ds := newDirSummary(nav)

	t.Run("Focus", func(t *testing.T) {
		ds.Focus(func(p tview.Primitive) {
			app.SetFocus(p)
		})
	})

	t.Run("selectionChanged", func(t *testing.T) {
		// Mock data to ensure we have rows
		entries := []os.DirEntry{
			mockDirEntry{name: "image1.png", isDir: false},
		}
		ds.SetDir(&DirContext{Path: "/test", children: entries})

		// Properly initialize nav.files and its rows to avoid panic in SetFilter
		nav.files.rows = NewFileRows(&DirContext{Path: "/test"})

		// We need at least one row in the table beyond the header
		if ds.extTable.GetRowCount() > 1 {
			ds.selectionChanged(1, 0) // Header is row 0
		}
	})

	t.Run("inputCapture", func(t *testing.T) {
		// Mock data with a group that has multiple extensions and a group that has one extension
		entries := []os.DirEntry{
			mockDirEntry{name: "image1.png", isDir: false},
			mockDirEntry{name: "image2.jpg", isDir: false},
			mockDirEntry{name: "script.go", isDir: false},
		}
		ds.SetDir(&DirContext{Path: "/test", children: entries})

		// Test Left
		eventLeft := tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
		assert.Nil(t, ds.inputCapture(eventLeft))

		// Test Down
		eventDown := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		// Row 0 is header group (Images or Code)
		ds.extTable.Select(0, 0)
		ds.inputCapture(eventDown)

		// Test Up
		eventUp := tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		ds.extTable.Select(1, 0)
		ds.inputCapture(eventUp)

		// Test other key
		eventOther := tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone)
		assert.Equal(t, eventOther, ds.inputCapture(eventOther))
	})

	t.Run("GetSizes", func(t *testing.T) {
		_ = ds.GetSizes()
	})
}
