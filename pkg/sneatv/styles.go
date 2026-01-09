package sneatv

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const DefaultFocusedBorderColor = tcell.ColorCornflowerBlue
const DefaultBlurBorderColor = tcell.ColorGray

func DefaultBorderWithPadding(box *tview.Box) {
	DefaultBorderWithoutPadding(box)
	box.SetBorderPadding(1, 1, 2, 2)
}

func DefaultBorderWithoutPadding(box *tview.Box) {
	box.SetBorder(true)
	box.SetBorderColor(DefaultBlurBorderColor)
	box.SetBorderAttributes(tcell.AttrDim)
	box.SetFocusFunc(func() {
		//box.SetBorderAttributes(tcell.AttrNone)
		box.SetBorderColor(DefaultFocusedBorderColor)
	})
	box.SetBlurFunc(func() {
		//box.SetBorderAttributes(tcell.AttrDim)
		box.SetBorderColor(DefaultBlurBorderColor)
	})
}

func SetPanelTitle(box *tview.Box, title string) {
	DefaultBorderWithPadding(box)
	box.SetTitle(title)
	box.SetTitleAlign(tview.AlignCenter)
	box.SetTitleColor(tview.Styles.TitleColor)
}
