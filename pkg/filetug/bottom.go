package filetug

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type bottom struct {
	*tview.TextView
}

type MenuItem struct {
	Title   string
	HotKeys []string
	Action  func()
}

func newBottom() *bottom {
	b := &bottom{
		TextView: tview.NewTextView().SetDynamicColors(true),
	}

	b.SetTextColor(tcell.ColorSlateGray)

	menuItems := []MenuItem{
		{
			Title:   "GoTo",
			HotKeys: []string{"G"},
			Action:  func() {},
		},
		{
			Title:   "F1 Help",
			HotKeys: []string{"F1"},
			Action:  func() {},
		},
		{
			Title:   "Preview",
			HotKeys: []string{"P"},
			Action:  func() {},
		},
	}

	const separator = "   "
	var sb strings.Builder
	for _, mi := range menuItems {
		title := mi.Title
		for _, key := range mi.HotKeys {
			title = strings.Replace(title, key, "[white]"+key+"[-]", 1)
		}
		sb.WriteString(title)
		sb.WriteString(separator)
	}
	b.SetText(sb.String()[:sb.Len()-len(separator)])
	return b
}
