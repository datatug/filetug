package sticky

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

type mockRecords struct {
	count int
}

func (m *mockRecords) RecordsCount() int {
	return m.count
}

func (m *mockRecords) GetCell(row, col int) *tview.TableCell {
	return tview.NewTableCell(fmt.Sprintf("R%dC%d", row, col))
}

func TestNewTable(t *testing.T) {
	columns := []Column{
		{Name: "Col1", Expansion: 1},
		{Name: "Col2", FixedWidth: 10},
	}
	table := NewTable(columns)
	assert.NotNil(t, table)
	assert.Equal(t, 2, table.GetColumnCount())

	// Check header cells
	cell0 := table.GetCell(0, 0)
	assert.Equal(t, "Col1", cell0.Text)
	cell1 := table.GetCell(0, 1)
	assert.Equal(t, "Col2", cell1.Text)
}

func TestTable_SetRecords(t *testing.T) {
	columns := []Column{{Name: "Col1"}}
	table := NewTable(columns)
	records := &mockRecords{count: 5}

	// We need to set a size for render to do something
	table.SetRect(0, 0, 100, 10)
	// Sticky table uses t.width which is set in DrawFunc
	// But it is also used in render() which is called by SetRecords.
	// In the current implementation, t.width might be 0 if Draw hasn't happened.

	table.SetRecords(records)

	// After SetRecords, render is called.
	// Since visibleRowsCount from GetRect (10) is > 0, it should render some rows.
	// Header is at row 0, records start at row 1.
	assert.Equal(t, 6, table.GetRowCount()) // 1 header + 5 records
	assert.Equal(t, "R0C0", table.GetCell(1, 0).Text)
}

func TestTable_ScrollToRow(t *testing.T) {
	columns := []Column{{Name: "Col1"}}
	table := NewTable(columns)
	records := &mockRecords{count: 100}
	table.SetRect(0, 0, 100, 10) // 10 rows total, 1 header -> 9 visible records
	table.SetRecords(records)

	// Initial state
	assert.Equal(t, 0, table.topRowIndex)

	// Scroll to row 20
	table.ScrollToRow(20)
	// topRowIndex should be 20 - 9 + 1 = 12
	assert.Equal(t, 12, table.topRowIndex)

	// Scroll back to row 5
	table.ScrollToRow(5)
	assert.Equal(t, 5, table.topRowIndex)

	// Scroll to row 10 (already visible since top=5, visible=9 -> 5..13)
	table.ScrollToRow(10)
	assert.Equal(t, 5, table.topRowIndex)
}

func TestTable_InputCapture(t *testing.T) {
	columns := []Column{{Name: "Col1"}}
	table := NewTable(columns)
	records := &mockRecords{count: 100}
	table.SetRect(0, 0, 100, 10)
	table.SetRecords(records)

	inputCapture := table.GetInputCapture()
	assert.NotNil(t, inputCapture)

	// Test KeyDown
	eventDown := tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
	inputCapture(eventDown)
	assert.Equal(t, 1, table.topRowIndex)

	// Test KeyUp
	eventUp := tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
	inputCapture(eventUp)
	assert.Equal(t, 0, table.topRowIndex)
}

func TestTable_Select(t *testing.T) {
	columns := []Column{{Name: "Col1"}}
	table := NewTable(columns)
	records := &mockRecords{count: 100}
	table.SetRect(0, 0, 100, 10)
	table.SetRecords(records)

	table.Select(20, 0)
	// Selecting row 20 (record 19) should trigger ScrollToRow(19)
	// topRowIndex should be 19 - 9 + 1 = 11
	assert.Equal(t, 11, table.topRowIndex)
}

func TestTable_Render_NoRecords(t *testing.T) {
	columns := []Column{{Name: "Col1"}}
	table := NewTable(columns)
	table.SetRect(0, 0, 100, 10)
	table.SetRecords(nil)
	assert.Equal(t, 1, table.GetRowCount()) // Just header
}
