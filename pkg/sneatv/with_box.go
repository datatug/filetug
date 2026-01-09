package sneatv

import (
	"github.com/rivo/tview"
)

type PrimitiveWithBox interface {
	tview.Primitive
	GetBox() *tview.Box
}

var _ PrimitiveWithBox = (*WithBoxType[tview.Primitive])(nil)

type WithBoxType[T tview.Primitive] struct {
	tview.Primitive
	Box *tview.Box
}

func (p WithBoxType[T]) GetBox() *tview.Box {
	return p.Box
}
func (p WithBoxType[T]) GetPrimitive() T {
	return p.Primitive.(T)
}

func WithDefaultBorders[T tview.Primitive](p T, box *tview.Box) WithBoxType[T] {
	DefaultBorderWithPadding(box)
	return WithBoxType[T]{
		Primitive: p,
		Box:       box,
	}
}

func WithBordersWithoutPadding[T tview.Primitive](p T, box *tview.Box) WithBoxType[T] {
	DefaultBorderWithoutPadding(box)
	return WithBoxType[T]{
		Primitive: p,
		Box:       box,
	}
}

func WithBoxWithoutBorder[T tview.Primitive](p T, box *tview.Box) WithBoxType[T] {
	return WithBoxType[T]{
		Primitive: p,
		Box:       box,
	}
}
