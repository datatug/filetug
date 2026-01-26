package viewers

import (
	"errors"
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/filetug/filetug/pkg/chroma2tcell"
	"github.com/filetug/filetug/pkg/files"
	"github.com/filetug/filetug/pkg/fsutils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var _ Previewer = (*TextPreviewer)(nil)

type TextPreviewer struct {
	*tview.TextView
}

func NewTextPreviewer() *TextPreviewer {
	return &TextPreviewer{
		TextView: tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(true).
			SetRegions(true).
			SetScrollable(true),
	}
}

func (p *TextPreviewer) Preview(entry files.EntryWithDirPath, data []byte, queueUpdateDraw func(func())) {
	go func() {
		if data == nil {
			var err error
			data, err = p.readFile(entry, 10*1024) // First 10KB
			if err != nil {
				return
			}
		}
		name := entry.Name()
		if lexer := lexers.Match(name); lexer == nil {
			queueUpdateDraw(func() {
				p.SetDynamicColors(false)
				p.TextView.SetText(string(data))
			})
		} else {
			colorized, err := chroma2tcell.Colorize(string(data), "dracula", lexer)
			queueUpdateDraw(func() {
				if err != nil {
					p.showError("Failed to format file: " + err.Error())
					return
				}
				p.TextView.Clear()
				p.TextView.SetDynamicColors(true)
				p.TextView.SetText(colorized)
				p.TextView.SetWrap(true)
			})
		}
	}()
}

func (p *TextPreviewer) Meta() tview.Primitive {
	return nil
}

func (p *TextPreviewer) Main() tview.Primitive {
	return p.TextView
}

func (p *TextPreviewer) readFile(entry files.EntryWithDirPath, max int) (data []byte, err error) {
	fullName := entry.FullName()
	data, err = fsutils.ReadFileData(fullName, max)
	if err != nil && !errors.Is(err, io.EOF) {
		p.showError(fmt.Sprintf("Failed to read file %s: %s", fullName, err.Error()))
		return
	}
	return
}

func (p *TextPreviewer) showError(text string) {
	p.SetDynamicColors(false)
	p.TextView.SetText(text)
	p.TextView.SetTextColor(tcell.ColorRed)
}
