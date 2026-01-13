package filetug

import (
	"os"
	"path"
	"time"

	"github.com/datatug/filetug/pkg/fsutils"
	"github.com/datatug/filetug/pkg/sticky"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var _ sticky.Records = (*FileRecords)(nil)

type FileRecords struct {
	NodePath string
	Entries  []os.DirEntry
	Infos    []os.FileInfo
	Err      error
}

func NewDirRecords(nodePath string, dirEntries []os.DirEntry) *FileRecords {
	return &FileRecords{
		NodePath: nodePath,
		Entries:  dirEntries,
		Infos:    make([]os.FileInfo, len(dirEntries)),
	}
}

func (r FileRecords) Count() int {
	if r.Err != nil {
		return 1
	}
	if len(r.Entries) == 0 {
		return 1
	}
	return len(r.Entries)
}

func (r FileRecords) GetCell(row, _ int, colName string) *tview.TableCell {
	if r.Err != nil {
		if colName == "Name" {
			return tview.NewTableCell(" 📁" + r.Err.Error()).SetTextColor(tcell.ColorOrangeRed)
		}
		return nil
	}
	if len(r.Entries) == 0 {
		if row == 0 && colName == "Name" {
			return tview.NewTableCell("[::i]No entries[::-]").SetTextColor(tcell.ColorGray)
		}
		return nil
	}
	dirEntry := r.Entries[row]
	var cell *tview.TableCell
	name := dirEntry.Name()
	if colName == "Name" {
		if dirEntry.IsDir() {
			cell = tview.NewTableCell(" 📁" + name)
		} else {
			cell = tview.NewTableCell(" 📄" + name)
		}
	} else {
		fi := r.Infos[row]
		if fi == nil {
			var err error
			fi, err = dirEntry.Info()
			if err != nil {
				return tview.NewTableCell(err.Error()).SetBackgroundColor(tcell.ColorRed)
			}
			r.Infos[row] = fi
		}

		switch colName {
		case "Size":
			var sizeText string
			if !dirEntry.IsDir() {
				size := fi.Size()
				sizeText = fsutils.GetSizeShortText(size)
			}
			cell = tview.NewTableCell(sizeText).SetAlign(tview.AlignRight)
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
	cell.SetReference(fsutils.ExpandHome(path.Join(r.NodePath, name)))
	return cell
}
