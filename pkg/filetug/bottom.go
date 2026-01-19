package filetug

import (
	"fmt"
	"strings"

	"github.com/datatug/filetug/pkg/ftui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type bottom struct {
	*tview.TextView
}

func newBottom() *bottom {
	b := &bottom{
		TextView: tview.NewTextView().SetDynamicColors(true),
	}

	b.SetTextColor(tcell.ColorSlateGray)

	menuItems := []ftui.MenuItem{
		{
			Title:   "F1Help",
			HotKeys: []string{"F1"},
			Action:  func() {},
		},
		{
			Title:   "Go",
			HotKeys: []string{"G"},
			Action:  func() {},
		},
		{
			Title:   "/root",
			HotKeys: []string{"/"},
			Action:  func() {},
		},
		{
			Title:   "~Home",
			HotKeys: []string{"H", "~"},
			Action:  func() {},
		},
		{
			Title:   "Favorites",
			HotKeys: []string{"F"},
			Action:  func() {},
		},
		{
			Title:   "Bookmarks",
			HotKeys: []string{"B"},
			Action:  func() {},
		},
		{
			Title:   "Lists",
			HotKeys: []string{"L"},
			Action:  func() {},
		},
		{
			Title:   "Previewer",
			HotKeys: []string{"P"},
			Action:  func() {},
		},
		{
			Title:   "Masks",
			HotKeys: []string{"M"},
			Action:  func() {},
		},
		{
			Title:   "Copy",
			HotKeys: []string{"F5", "C"},
			Action:  func() {},
		},
		{
			Title:   "Rename",
			HotKeys: []string{"F6", "R"},
			Action:  func() {},
		},
		{
			Title:   "Delete",
			HotKeys: []string{"F8", "D"},
			Action:  func() {},
		},
		{
			Title:   "View",
			HotKeys: []string{"V"},
			Action:  func() {},
		},
		{
			Title:   "Edit",
			HotKeys: []string{"E"},
			Action:  func() {},
		},
		{
			Title:   "Exit",
			HotKeys: []string{"x"},
			Action:  func() {},
		},
	}

	const separator = "┊"
	var sb strings.Builder
	for _, mi := range menuItems {
		title := mi.Title
		for _, key := range mi.HotKeys {
			title = strings.Replace(title, key, fmt.Sprintf("[%s]%s[-]", ftui.CurrentTheme.HotkeyColor, key), 1)
		}
		sb.WriteString(title)
		sb.WriteString(separator)
	}
	b.SetText(sb.String()[:sb.Len()-len(separator)])
	return b
}
