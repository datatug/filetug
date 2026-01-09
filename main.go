package main

import (
	"os"

	"github.com/datatug/filetug/pkg/filetug"
	"github.com/rivo/tview"
)

func main() {
	app := newApp()
	run(app)
}

var setupApp = filetug.SetupApp

var newApp = func() *tview.Application {
	app := tview.NewApplication()
	setupApp(app)
	return app
}

type application interface{ Run() error }

var run = func(app application) {
	if err := app.Run(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
	}
}
