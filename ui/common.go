package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := image.Point{X: cs.Min.X, Y: cs.Min.Y}
	st := op.Save(gtx.Ops)
	track := image.Rectangle{
		Max: d,
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	st.Load()

	return layout.Dimensions{Size: d}
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func endToEndRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.W.Layout(gtx, leftWidget)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.E.Layout(gtx, rightWidget)
		}),
	)
}

type (
	D = layout.Dimensions
	C = layout.Context
)

var (
	contentAdd    *widget.Icon
	contentRemove *widget.Icon
	contentSave   *widget.Icon
)

func init() {
	contentAdd, _ = widget.NewIcon(icons.ContentAdd)
	contentRemove, _ = widget.NewIcon(icons.ContentRemove)
	contentSave, _ = widget.NewIcon(icons.ContentSave)
}
