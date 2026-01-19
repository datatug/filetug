package masks

import (
	"github.com/datatug/filetug/pkg/ftui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Selected struct {
	Mask     *Mask
	Patterns []*Pattern
}

type Panel struct {
	*tview.Table
	boxed *ftui.Boxed
	masks []Mask
}

func (p *Panel) Draw(screen tcell.Screen) {
	p.boxed.Draw(screen)
}

func (p *Panel) Focus(delegate func(p tview.Primitive)) {
	p.Table.Focus(delegate)
}

func NewPanel() *Panel {
	p := new(Panel)
	p.masks = createBuiltInMasks()

	p.Table = tview.NewTable()
	p.SetTitle("Masks")
	p.SetSelectable(true, true)
	p.SetFixed(1, 1)

	p.boxed = ftui.NewBoxed(p.Table,
		ftui.WithLeftBorder(0, -1),
	)

	p.SetCell(0, 0, tview.NewTableCell("Mask").SetExpansion(1))
	p.SetCell(0, 1, tview.NewTableCell("CurrDir").SetAlign(tview.AlignRight))
	p.SetCell(0, 2, tview.NewTableCell("SubDirs").SetAlign(tview.AlignRight))

	for i, m := range p.masks {

		nameCell := tview.NewTableCell(m.Name)
		nameCell.SetExpansion(1)
		p.SetCell(i+1, 0, nameCell)

		currDirCell := tview.NewTableCell("...")
		currDirCell.SetAlign(tview.AlignRight)
		p.SetCell(i+1, 1, currDirCell)

		subDirsCell := tview.NewTableCell("...")
		subDirsCell.SetAlign(tview.AlignRight)
		subDirsCell.SetTextColor(tcell.ColorGray)
		p.SetCell(i+1, 2, subDirsCell)
	}

	return p
}
