package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type Dialog struct {
	modal *component.ModalLayer
	btn   *widget.Clickable
	pwd   *widget.Editor
}

func (d *Dialog) Widget(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
	for d.btn.Clicked() {
		d.modal.Disappear(gtx.Now)
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		clip.Rect{Max: image.Pt(150, 150)}.Push(gtx.Ops)
		paint.ColorOp{Color: th.Bg}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Body2(th, "Input Password").Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Editor(th, d.pwd, "please admin password").Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, d.btn, "cancel").Layout(gtx)
			}),
		)

	})
}
