package sneatv

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ButtonWithShortcut struct {
	*tview.Button
	shortcut              rune
	shortcutStyle         tcell.Style
	shortcutActivateStyle tcell.Style
}

func NewButtonWithShortcut(label string, shortcut rune) *ButtonWithShortcut {
	/*
		mainTextStyle:      ,
		secondaryTextStyle: tcell.StyleDefault.Foreground(Styles.TertiaryTextColor).Background(Styles.PrimitiveBackgroundColor),

	*/
	return &ButtonWithShortcut{
		shortcut: shortcut,
		shortcutStyle: tcell.StyleDefault.
			Foreground(tview.Styles.SecondaryTextColor).
			Background(tview.Styles.PrimitiveBackgroundColor),
		shortcutActivateStyle: tcell.StyleDefault.
			Background(tcell.ColorYellow).
			Foreground(tcell.ColorWhite),
		Button: tview.NewButton(label),
	}
}

func (b *ButtonWithShortcut) SetShortcutStyle(style tcell.Style) *ButtonWithShortcut {
	b.shortcutStyle = style
	return b
}

// Draw draws this primitive onto the screen.
func (b *ButtonWithShortcut) Draw(screen tcell.Screen) {
	// Draw the box.
	style := tcell.StyleDefault.Background(tview.Styles.PrimitiveBackgroundColor).Foreground(tview.Styles.PrimaryTextColor)
	disabledStyle := tcell.StyleDefault.Background(tview.Styles.PrimitiveBackgroundColor).Foreground(tcell.ColorGray)
	shortcutStyle := b.shortcutStyle

	if b.IsDisabled() {
		style = disabledStyle
	}
	if b.HasFocus() && !b.IsDisabled() {
		style = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack)
		shortcutStyle = b.shortcutActivateStyle
	}
	// Avoid deprecated Style.Decompose(): derive background color from state
	bgColor := tview.Styles.PrimitiveBackgroundColor
	if b.HasFocus() && !b.IsDisabled() {
		bgColor = tcell.ColorYellow
	}
	b.SetBackgroundColor(bgColor)
	b.DrawForSubclass(screen, b)

	// Draw label with shortcut
	x, y, width, height := b.GetInnerRect()
	if width > 0 && height > 0 {
		y = y + height/2

		// Format shortcut and label parts separately
		shortcutText := fmt.Sprintf("(%c)", b.shortcut)
		labelText := fmt.Sprintf(" %s", b.GetLabel())

		// Calculate starting position for centered text
		totalWidth := len(shortcutText) + len(labelText)
		startX := x + (width-totalWidth)/2

		if startX < x {
			startX = x
		}

		// Render shortcut part with shortcutStyle
		currentX := startX

		if shortcutStyle == (tcell.Style{}) {
			shortcutStyle = style // fallback to button style if shortcutStyle is not set
		}

		for _, ch := range shortcutText {
			if currentX < x+width {
				screen.SetContent(currentX, y, ch, nil, shortcutStyle)
				currentX++
			}
		}

		// Render label part with regular button style
		for _, ch := range labelText {
			if currentX < x+width {
				screen.SetContent(currentX, y, ch, nil, style)
				currentX++
			}
		}
	}
}
