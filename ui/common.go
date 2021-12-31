package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

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
