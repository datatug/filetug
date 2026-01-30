package sneatv

import "github.com/rivo/tview"

// TabsApp defines what methods of a tview.Application are used by tabs.
// User navigator.NewApp with options or navigator.MockApp for tests
type TabsApp interface {
	QueueUpdateDraw(f func())
	SetFocus(p tview.Primitive)
}
