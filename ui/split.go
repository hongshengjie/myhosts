package ui

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Split struct {
	// Ratio keeps the current layout.
	// 0 is center, -1 completely to the left, 1 completely to the right.
	Ratio float32
	// Bar is the width for resizing the layout
	Bar unit.Value

	drag   bool
	dragID pointer.ID
	dragX  float32
}

var defaultBarWidth = unit.Dp(5)

func (s *Split) Layout(gtx layout.Context, th *material.Theme, left, right layout.Widget) layout.Dimensions {
	bar := gtx.Px(s.Bar)
	if bar <= 1 {
		bar = gtx.Px(defaultBarWidth)
	}

	proportion := (s.Ratio + 1) / 2
	leftsize := int(proportion*float32(gtx.Constraints.Max.X) - float32(bar))

	rightoffset := leftsize + bar
	rightsize := gtx.Constraints.Max.X - rightoffset

	{ // handle input
		// Avoid affecting the input tree with pointer events.
		stack := op.Save(gtx.Ops)

		for _, ev := range gtx.Events(s) {
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			switch e.Type {
			case pointer.Press:
				if s.drag {
					break
				}

				s.dragID = e.PointerID
				s.dragX = e.Position.X

			case pointer.Drag:
				if s.dragID != e.PointerID {
					break
				}

				deltaX := e.Position.X - s.dragX
				s.dragX = e.Position.X

				deltaRatio := deltaX * 2 / float32(gtx.Constraints.Max.X)
				s.Ratio += deltaRatio

			case pointer.Release:
				fallthrough
			case pointer.Cancel:
				s.drag = false
			}
		}

		// register for input
		barRect := image.Rect(leftsize, gtx.Constraints.Min.Y, rightoffset, gtx.Constraints.Max.X)

		pointer.Rect(barRect).Add(gtx.Ops)
		pointer.InputOp{Tag: s,
			Types: pointer.Press | pointer.Drag | pointer.Release,
			Grab:  s.drag,
		}.Add(gtx.Ops)
		r := clip.Rect(barRect).Op()
		paint.FillShape(gtx.Ops, th.ContrastBg, r)
		stack.Load()
	}

	{
		stack := op.Save(gtx.Ops)

		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftsize, gtx.Constraints.Max.Y))
		left(gtx)

		stack.Load()
	}

	{
		stack := op.Save(gtx.Ops)

		op.Offset(f32.Pt(float32(rightoffset), 0)).Add(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		right(gtx)

		stack.Load()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
