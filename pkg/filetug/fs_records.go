package filetug

import (
	"os"
	"path"
	"time"

	"github.com/datatug/filetug/pkg/fsutils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//var _ sticky.Records = (*FileRows)(nil)

var _ tview.TableContent = (*FileRows)(nil)

type FileRows struct {
	tview.TableContentReadOnly
	NodePath string
	Entries  []os.DirEntry
	Infos    []os.FileInfo
	//selected string
	Err error
}

//func (r *FileRows) SetSelected(row int) {
//	if row == 0 {
//		r.selected = ""
//	}
//	r.selected = r.Entries[row-1].Name()
//}

func (r *FileRows) GetRowCount() int {
	return len(r.Entries) + 1
}

func (r *FileRows) GetColumnCount() int {
	return 3
}

func NewFileRows(nodePath string, dirEntries []os.DirEntry) *FileRows {
	return &FileRows{
		NodePath: nodePath,
		Entries:  dirEntries,
		Infos:    make([]os.FileInfo, len(dirEntries)),
	}
}

const (
	nameColIndex     = 0
	sizeColIndex     = 1
	modifiedColIndex = 2
)

func (r *FileRows) GetCell(row, col int) *tview.TableCell {
	if row < 0 {
		return nil
	}
	if row == 0 {
		th := func(text string) *tview.TableCell {
			return tview.NewTableCell(text)
		}
		switch col {
		case nameColIndex:
			return th("Name").SetExpansion(1)
		case sizeColIndex:
			return th("Size")
		case modifiedColIndex:
			return th("Modified")
		default:
			return nil
		}
	}
	if r.Err != nil {
		if col == nameColIndex {
			return tview.NewTableCell(" 📁" + r.Err.Error()).SetTextColor(tcell.ColorOrangeRed)
		}
		return nil
	}
	if len(r.Entries) == 0 {
		if col == nameColIndex {
			return tview.NewTableCell("[::i]No entries[::-]").SetTextColor(tcell.ColorGray)
		}
		return nil
	}
	i := row - 1
	dirEntry := r.Entries[i]
	var cell *tview.TableCell
	name := dirEntry.Name()
	if col == nameColIndex {
		if dirEntry.IsDir() {
			cell = tview.NewTableCell(" 📁" + name)
		} else {
			cell = tview.NewTableCell(" 📄" + name)
		}
	} else {
		fi := r.Infos[i]
		if fi == nil {
			var err error
			fi, err = dirEntry.Info()
			if err != nil {
				return tview.NewTableCell(err.Error()).SetBackgroundColor(tcell.ColorRed)
			}
			r.Infos[i] = fi
		}

		switch col {
		case sizeColIndex:
			var sizeText string
			if !dirEntry.IsDir() {
				size := fi.Size()
				sizeText = fsutils.GetSizeShortText(size)
			}
			cell = tview.NewTableCell(sizeText).SetAlign(tview.AlignRight)
		case modifiedColIndex:
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
	cell.SetReference(fsutils.ExpandHome(path.Join(r.NodePath, name)))
	return cell
}
