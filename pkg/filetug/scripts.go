package filetug

import (
	"github.com/filetug/filetug/pkg/sneatv"
	"github.com/rivo/tview"
)

type scriptsPanel struct {
	*sneatv.Boxed
	list *tview.List
}

func (nav *Navigator) showScriptsPanel() {
	list := tview.NewList()
	list.AddItem("Nested Dirs Generator", "", '1', func() {
		showNestedDirsGenerator(nav)
	})
	p := &scriptsPanel{
		Boxed: sneatv.NewBoxed(list),
		list:  list,
	}
	nav.right.SetContent(p)
	nav.app.SetFocus(p)
}

func showNestedDirsGenerator(nav *Navigator) {
	currentBrowser := nav.getCurrentBrowser()
	p := newNestedDirsGeneratorPanel(nav, currentBrowser)
	nav.right.SetContent(p)
	nav.app.SetFocus(p)
}
