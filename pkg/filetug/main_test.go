package filetug

import (
	"testing"

	"github.com/rivo/tview"
)

func TestSetupApp(t *testing.T) {
	app := tview.NewApplication()
	SetupApp(app)
}
