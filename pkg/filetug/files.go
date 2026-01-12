package filetug

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/datatug/filetug/pkg/fsutils"
	"github.com/datatug/filetug/pkg/sticky"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var _ sticky.Records = (*fsRecords)(nil)

type fsRecords struct {
	nodePath   string
	dirEntries []os.DirEntry
	infos      []os.FileInfo
}

func NewDirRecords(nodePath string, dirEntries []os.DirEntry) sticky.Records {
	return &fsRecords{
		nodePath:   nodePath,
		dirEntries: dirEntries,
		infos:      make([]os.FileInfo, len(dirEntries)),
	}
}

func (r fsRecords) Count() int {
	return len(r.dirEntries)
}

func (r fsRecords) GetCell(row, _ int, colName string) *tview.TableCell {
	dirEntry := r.dirEntries[row]
	var cell *tview.TableCell
	name := dirEntry.Name()
	if colName == "Name" {
		if dirEntry.IsDir() {
			cell = tview.NewTableCell(" 📁" + name)
		} else {
			cell = tview.NewTableCell(" 📄" + name)
		}
	} else {
		fi := r.infos[row]
		if fi == nil {
			var err error
			fi, err = dirEntry.Info()
			if err != nil {
				return tview.NewTableCell(err.Error()).SetBackgroundColor(tcell.ColorRed)
			}
			r.infos[row] = fi
		}

		switch colName {
		case "Size":
			cell = tview.NewTableCell(strconv.FormatInt(fi.Size(), 10)).SetAlign(tview.AlignRight)
		case "Modified":
			var s string
			if modTime := fi.ModTime(); fi.ModTime().After(time.Now().Add(24 * time.Hour)) {
				s = modTime.Format("15:04:05")
			} else {
				s = modTime.Format("2006-01-02")
			}
			cell = tview.NewTableCell(s)
		default:
			return nil
		}
	}
	color := GetColorByFileExt(name)
	cell.SetTextColor(color)
	cell.SetReference(fsutils.ExpandHome(path.Join(r.nodePath, name)))
	return cell
}

type files struct {
	*sticky.Table
	nav   *Navigator
	boxed *boxed
}

func (f *files) Draw(screen tcell.Screen) {
	f.boxed.Draw(screen)
}

func (f *files) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	table := f.Table
	if string(event.Rune()) == " " {
		row, _ := table.GetSelection()
		cell := table.GetCell(row, 0)

		if strings.HasPrefix(cell.Text, " ") {
			cell.SetText("✓" + strings.TrimPrefix(cell.Text, " "))
		} else {
			cell.SetText(" " + strings.TrimPrefix(cell.Text, "✓"))
		}
		return nil
	}
	switch event.Key() {
	case tcell.KeyLeft:
		f.nav.app.SetFocus(f.nav.dirsTree)
		return nil
	case tcell.KeyRight:
		f.nav.app.SetFocus(f.nav.previewer)
		return nil
	case tcell.KeyUp:
		row, _ := table.GetSelection()
		if row == 0 {
			if f.nav.o.moveFocusUp != nil {
				f.nav.o.moveFocusUp(table)
				return nil
			}
		}
		return event
	default:
		return event
	}
}

func newFiles(nav *Navigator) *files {
	table := sticky.NewTable([]sticky.Column{
		{
			Name:      "Name",
			Expansion: 1,
			MinWidth:  20,
		},
		{
			Name:       "Size",
			FixedWidth: 6,
		},
		{
			Name:       "Modified",
			FixedWidth: 10,
		},
	})
	f := &files{
		nav:   nav,
		Table: table,
		boxed: newBoxed(
			table,
			WithLeftBorder(0, -1),
			WithRightBorder(0, +1),
		),
	}
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetInputCapture(f.inputCapture)
	table.SetFocusFunc(func() {
		nav.activeCol = 1
	})
	nav.filesFocusFunc = func() {
		nav.activeCol = 1
	}

	table.SetSelectionChangedFunc(f.selectionChanged)
	nav.filesSelectionChangedFunc = f.selectionChangedNavFunc
	return f
}

// selectionChangedNavFunc: TODO: is it a duplicate of selectionChangedNavFunc?
func (f *files) selectionChangedNavFunc(row, _ int) {
	if row == 0 {
		f.nav.previewer.textView.SetText("Selected dir: " + f.nav.currentDir)
		f.nav.previewer.textView.SetTextColor(tcell.ColorWhiteSmoke)
		return
	}
	cell := f.GetCell(row, 0)
	name := cell.Text[1:]
	fullName := filepath.Join(f.nav.currentDir, name)
	f.nav.previewer.PreviewFile(name, fullName)
}

// selectionChanged: TODO: is it a duplicate of selectionChangedNavFunc?
func (f *files) selectionChanged(row, _ int) {
	if row == 0 {
		f.nav.previewer.textView.SetText("Selected dir: " + f.nav.currentDir)
		f.nav.previewer.textView.SetTextColor(tcell.ColorWhiteSmoke)
		return
	}
	cell := f.GetCell(row, 0)
	ref := cell.GetReference()
	if ref == nil {
		f.nav.previewer.SetText("cell has no reference")
		return
	}
	fullName := ref.(string)
	stat, err := os.Stat(fullName)
	if err != nil {
		f.nav.previewer.SetErr(err)
		return
	}
	if stat.IsDir() {
		f.nav.previewer.SetText("Directory: " + fullName)
		return
	}
	f.nav.previewer.PreviewFile("", fullName)
}
