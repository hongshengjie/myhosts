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

type PasswordDialog struct {
	modal     *component.ModalLayer
	cancelbtn *widget.Clickable
	pwd       *widget.Editor
}

func (d *PasswordDialog) Widget(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
	for d.cancelbtn.Clicked() {
		d.modal.Disappear(gtx.Now)
	}
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		clip.Rect{Max: image.Pt(150, 150)}.Push(gtx.Ops)
		paint.ColorOp{Color: th.Bg}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.H6(th, "Password").Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Editor(th, d.pwd, "input password").Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, d.cancelbtn, "cancel").Layout(gtx)
			}),
		)

	})
}

type TitleDialog struct {
	modal      *component.ModalLayer
	cancelbtn  *widget.Clickable
	titleInput *widget.Editor
}

func (d *TitleDialog) Widget(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
	for d.cancelbtn.Clicked() {
		d.modal.Disappear(gtx.Now)
	}
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		clip.Rect{Max: image.Pt(150, 150)}.Push(gtx.Ops)
		paint.ColorOp{Color: th.Bg}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.H6(th, "Title").Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Editor(th, d.titleInput, "input title").Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, d.cancelbtn, "cancel").Layout(gtx)
			}),
		)

	})
}
