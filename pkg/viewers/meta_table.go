package viewers

import "github.com/rivo/tview"

func NewMetaTable() *MetaTable {
	return &MetaTable{
		Table: tview.NewTable(),
	}
}

type MetaTable struct {
	*tview.Table
	meta *Meta
}

func (mt *MetaTable) SetMeta(meta *Meta) {
	mt.meta = meta
	mt.Clear()

	row := 0

	for _, group := range meta.Groups {
		{
			groupCell := tview.NewTableCell(group.Title)
			mt.SetCell(row, 0, groupCell)
			row++
		}
		for _, record := range group.Records {
			{ // Title cell
				titleCell := tview.NewTableCell("  " + record.Title)
				titleCell.Align = tview.AlignRight
				mt.SetCell(row, 0, titleCell)
			}
			{ // Value cell
				valueCell := tview.NewTableCell(record.Value)
				switch record.ValueAlign {
				case AlignRight:
					valueCell.Align = tview.AlignRight
				case AlignLeft: // Do nothing
				}
				mt.SetCell(row, 1, valueCell)
			}
			row++
		}
	}
}
