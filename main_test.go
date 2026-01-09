package main

import (
	"testing"

	"github.com/datatug/filetug/pkg/filetug"
	"github.com/rivo/tview"
)

func TestMainRoot(t *testing.T) {
	app := tview.NewApplication()
	go func() {
		filetug.Main(app)
	}()
	app.Stop()
}
