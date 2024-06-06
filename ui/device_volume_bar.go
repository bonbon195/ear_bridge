package ui

import (
	"log"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/widget"
)

func progressBarButton(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	if button.Clicked(gtx) {
		presses := button.History()
		pos := presses[len(presses)-1].Position.X
		posRatio := float32(pos) / float32(gtx.Constraints.Max.X)
		defer func() {
			recover()
		}()
		log.Println(posRatio)
	}

	return button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		return layout.UniformInset(0).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return w(gtx)
		})
	})
}
