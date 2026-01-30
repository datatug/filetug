package filetug

import (
	"github.com/filetug/filetug/pkg/filetug/navigator"
)

func SetupApp(app navigator.App) {
	app.EnableMouse(true)
	nav := NewNavigator(app)
	app.SetRoot(nav, true)
}
