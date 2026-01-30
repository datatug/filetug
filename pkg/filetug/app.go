package filetug

import "github.com/rivo/tview"

type ftApp struct {
	*tview.Application
}

func (a ftApp) QueueUpdateDraw(f func()) {
	if a.Application != nil {
		_ = a.Application.QueueUpdateDraw(f)
	} else {
		f()
	}
}

func (a ftApp) SetFocus(p tview.Primitive) {
	if a.Application != nil {
		_ = a.Application.SetFocus(p)
	}
}
