package sneatv

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNewTabs(t *testing.T) {
	tabs := NewTabs(nil, UnderlineTabsStyle, WithLabel("Tabs:"))
	assert.NotNil(t, tabs)
	assert.Equal(t, "Tabs:", tabs.label)
}

func TestTabs_AddAndSwitch(t *testing.T) {
	tabs := NewTabs(nil, UnderlineTabsStyle)
	tab1 := &Tab{ID: "1", Title: "Tab 1", Primitive: tview.NewBox()}
	tab2 := &Tab{ID: "2", Title: "Tab 2", Primitive: tview.NewBox()}
	tabs.AddTabs(tab1, tab2)

	assert.Equal(t, 2, len(tabs.tabs))
	assert.Equal(t, 0, tabs.active)

	tabs.SwitchTo(1)
	assert.Equal(t, 1, tabs.active)

	tabs.SwitchTo(5) // out of bounds
	assert.Equal(t, 1, tabs.active)
}

func TestTabs_Navigation(t *testing.T) {
	tabs := NewTabs(nil, UnderlineTabsStyle,
		FocusLeft(func(current tview.Primitive) {}),
		FocusUp(func(current tview.Primitive) {}),
		FocusDown(func(current tview.Primitive) {}),
	)
	tabs.AddTabs(
		&Tab{ID: "1", Title: "T1", Primitive: tview.NewBox()},
		&Tab{ID: "2", Title: "T2", Primitive: tview.NewBox()},
	)

	// Right
	event := tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
	res := tabs.handleInput(event)
	assert.Nil(t, res)
	assert.Equal(t, 1, tabs.active)

	// Left
	event = tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
	res = tabs.handleInput(event)
	assert.Nil(t, res)
	assert.Equal(t, 0, tabs.active)

	// FocusLeft
	leftCalled := false
	tabs.focusLeft = func(current tview.Primitive) { leftCalled = true }
	tabs.handleInput(event)
	assert.True(t, leftCalled)

	// Alt+1
	event = tcell.NewEventKey(tcell.KeyRune, '1', tcell.ModAlt)
	tabs.handleInput(event)
	assert.Equal(t, 0, tabs.active)
}

func TestTabs_FocusOptions(t *testing.T) {
	downCalled := false
	tabs := NewTabs(nil, UnderlineTabsStyle, FocusDown(func(current tview.Primitive) {
		downCalled = true
	}))
	tabs.AddTabs(&Tab{ID: "1", Title: "T1", Primitive: tview.NewBox()})

	event := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	tabs.handleInput(event)
	assert.True(t, downCalled)
}
