package navigator

import "github.com/rivo/tview"

type UpdateDrawQueuer func(f func())

type Focuser func(primitive tview.Primitive)

type RootSetter func(root tview.Primitive, fullscreen bool)
