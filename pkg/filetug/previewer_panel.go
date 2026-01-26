package filetug

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/filetug/filetug/pkg/files"
	"github.com/filetug/filetug/pkg/sneatv"
	"github.com/filetug/filetug/pkg/viewers"
	"github.com/filetug/filetug/pkg/viewers/imageviewer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type previewerPanel struct {
	*sneatv.Boxed
	rows      *tview.Flex
	nav       *Navigator
	attrsRow  *tview.Flex
	fsAttrs   *tview.Table
	separator *tview.TextView
	previewer viewers.Previewer
	textView  *tview.TextView
}

func newPreviewerPanel(nav *Navigator) *previewerPanel {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	p := previewerPanel{
		Boxed: sneatv.NewBoxed(
			flex,
			sneatv.WithLeftBorder(0, -1),
		),
		rows:      flex,
		attrsRow:  tview.NewFlex().SetDirection(tview.FlexRow),
		fsAttrs:   tview.NewTable(),
		separator: tview.NewTextView().SetText(strings.Repeat("─", 20)).SetTextColor(tcell.ColorGray),
		textView:  tview.NewTextView(),
		nav:       nav,
	}

	p.textView.SetWrap(false)
	p.textView.SetDynamicColors(true)
	p.textView.SetText("To be implemented.")
	p.textView.SetFocusFunc(func() {
		nav.activeCol = 2
	})

	p.attrsRow.AddItem(p.fsAttrs, 0, 1, false)

	p.rows.AddItem(p.attrsRow, 2, 0, false)
	p.rows.AddItem(p.separator, 1, 0, false)
	//p.rows.AddItem(p.textView, 0, 1, false)

	p.rows.SetFocusFunc(func() {
		nav.activeCol = 2
		p.rows.SetBorderColor(sneatv.CurrentTheme.FocusedBorderColor)
	})
	nav.previewerFocusFunc = func() {
		nav.activeCol = 2
		p.rows.SetBorderColor(sneatv.CurrentTheme.FocusedBorderColor)
	}
	p.rows.SetBlurFunc(func() {
		p.rows.SetBorderColor(sneatv.CurrentTheme.BlurredBorderColor)
	})
	nav.previewerBlurFunc = func() {
		p.rows.SetBorderColor(sneatv.CurrentTheme.BlurredBorderColor)
	}

	p.rows.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			nav.setAppFocus(nav.files)
			return nil
		case tcell.KeyUp:
			nav.o.moveFocusUp(p.textView)
			return nil
		default:
			return event
		}
	})

	return &p
}

func (p *previewerPanel) setPreviewer(previewer viewers.Previewer) {
	if p.previewer != nil {
		if meta := p.previewer.Meta(); meta != nil {
			p.attrsRow.RemoveItem(meta)
		}
		if main := p.previewer.Main(); main != nil {
			p.rows.RemoveItem(main)
		}
	}
	p.previewer = previewer
	if previewer != nil {
		if meta := previewer.Meta(); meta != nil {
			p.attrsRow.AddItem(meta, 0, 1, false)
		}
		if main := previewer.Main(); main != nil {
			p.rows.AddItem(main, 0, 1, false)
		}
	}

}

func (p *previewerPanel) SetErr(err error) {
	p.textView.Clear()
	p.textView.SetDynamicColors(true)
	p.textView.SetText(err.Error())
	p.textView.SetTextColor(tcell.ColorRed)
}

func (p *previewerPanel) SetText(text string) {
	p.textView.Clear()
	p.textView.SetDynamicColors(true)
	p.textView.SetText(text)
	p.textView.SetTextColor(tcell.ColorWhiteSmoke)
}

func (p *previewerPanel) PreviewEntry(entry files.EntryWithDirPath) {
	name := entry.Name()
	fullName := entry.FullName()
	if name == "" {
		_, name = path.Split(fullName)
	}
	p.SetTitle(name)
	switch name {
	case ".DS_Store":
		p.previewer = viewers.NewDsstorePreviewer()
	default:
		ext := strings.ToLower(filepath.Ext(name))
		switch ext {
		case ".json":
			p.setPreviewer(viewers.NewJsonPreviewer())
		case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".riff", ".tiff", ".vp8", ".webp":
			p.setPreviewer(imageviewer.NewImagePreviewer())
			return
		default:
			p.setPreviewer(viewers.NewTextPreviewer())
		}
	}
	p.previewer.Preview(entry, nil, p.nav.queueUpdateDraw)
}
