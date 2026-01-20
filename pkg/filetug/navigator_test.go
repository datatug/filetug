package filetug

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/datatug/filetug/pkg/gitutils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestOnMoveFocusUp(t *testing.T) {
	var s tview.Primitive
	f := func(source tview.Primitive) {
		s = source
	}
	o := OnMoveFocusUp(f)
	var options navigatorOptions
	o(&options)
	assert.Equal(t, f, options.moveFocusUp)

	textView := tview.NewTextView()
	options.moveFocusUp(textView)
	assert.Equal[tview.Primitive](t, textView, s)
}

func TestNavigator(t *testing.T) {
	app := tview.NewApplication()
	nav := NewNavigator(app, OnMoveFocusUp(func(source tview.Primitive) {}))
	if nav == nil {
		t.Fatal("nav is nil")
	}

	t.Run("SetFocus", func(t *testing.T) {
		nav.SetFocus()
	})

	t.Run("NavigatorInputCapture", func(t *testing.T) {
		altKey := func(r rune) *tcell.EventKey {
			return tcell.NewEventKey(tcell.KeyRune, r, tcell.ModAlt)
		}
		nav.GetInputCapture()(altKey('0'))
		nav.GetInputCapture()(altKey('+'))
		nav.GetInputCapture()(altKey('-'))

		nav.activeCol = 1
		nav.GetInputCapture()(altKey('+'))
		nav.GetInputCapture()(altKey('-'))

		nav.activeCol = 2
		nav.GetInputCapture()(altKey('+'))
		nav.GetInputCapture()(altKey('-'))

		nav.activeCol = -1
		nav.GetInputCapture()(altKey('+'))
		nav.GetInputCapture()(altKey('-'))

		nav.GetInputCapture()(altKey('f'))
		nav.GetInputCapture()(altKey('m'))
		nav.GetInputCapture()(altKey('r'))
		nav.GetInputCapture()(altKey('h'))
		nav.GetInputCapture()(altKey('?'))
		nav.GetInputCapture()(altKey('z')) // unknown alt key

		nav.GetInputCapture()(tcell.NewEventKey(tcell.KeyF1, 0, tcell.ModNone))
		nav.GetInputCapture()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))

		// Test moveFocusUp in navigator
		nav.o.moveFocusUp(nav.files)

		// Test Ctrl modifier
		nav.GetInputCapture()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModCtrl))
		assert.True(t, nav.bottom.isCtrl)
	})
}

func TestNavigator_GitStatus(t *testing.T) {
	app := tview.NewApplication()
	nav := NewNavigator(app)
	if nav == nil {
		t.Fatal("nav is nil")
	}
	node := tview.NewTreeNode("test")

	drawUpdatesCount := 0
	oldQueueUpdateDraw := nav.queueUpdateDraw
	defer func() {
		nav.queueUpdateDraw = oldQueueUpdateDraw
	}()
	nav.queueUpdateDraw = func(f func()) *tview.Application {
		drawUpdatesCount++
		return app
	}

	// Use background context for tests
	ctx := context.Background()

	// 1. Not in cache, git status returns nil
	nav.updateGitStatus(ctx, "/non-existent", node, "prefix: ")

	// 2. In cache
	nav.gitStatusCache["/cached"] = &gitutils.RepoStatus{Branch: "main"}
	nav.updateGitStatus(ctx, "/cached", node, "prefix: ")

	// 3. Cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	nav.updateGitStatus(cancelledCtx, "/any", node, "prefix: ")

	time.Sleep(100 * time.Millisecond)
}

func TestNavigator_goDir(t *testing.T) {
	saveCurrentDir = func(string, string) {}
	app := tview.NewApplication()
	nav := NewNavigator(app, OnMoveFocusUp(func(source tview.Primitive) {}))

	t.Run("goDir_Success", func(t *testing.T) {
		nav.goDir(".")
	})

	t.Run("goDir_NonExistent", func(t *testing.T) {
		nav.goDir("/non-existent-Path-12345")
	})

	t.Run("Extra", func(t *testing.T) {
		nav.SetFocusToContainer(0)
		nav.SetFocusToContainer(1)
		nav.SetFocusToContainer(2)
		nav.showMasks()
	})

	t.Run("onDataLoaded_showNodeError", func(t *testing.T) {
		node := tview.NewTreeNode("test").SetReference("/test")
		dirContext := &DirContext{
			Path:     "/test",
			children: []os.DirEntry{mockDirEntry{name: "file.txt", isDir: false}},
		}

		nav.onDataLoaded(node, dirContext, true)
		nav.onDataLoaded(node, dirContext, false)

		err := errors.New("test error")
		nav.showNodeError(node, err)
		nav.showNodeError(nil, err)
	})
}
