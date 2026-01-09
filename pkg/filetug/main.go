package filetug

import (
	"github.com/rivo/tview"
)

func SetupApp(app *tview.Application) {
	app.EnableMouse(true)
	app.SetRoot(NewNavigator(app), true)
}
