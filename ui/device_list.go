package ui

import (
	"image/color"
	"slices"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/bonbon195/ear_bridge/audio"
	"github.com/bonbon195/ear_bridge/ui/style"
	"github.com/gen2brain/malgo"
	"github.com/nuttech/bell/v2"
)

type DeviceList struct {
	Devices    []DeviceListItem
	LayoutList *widget.List
}

type DeviceListItem struct {
	Id     malgo.DeviceID
	Name   string
	Chosen bool
	Button *widget.Clickable
}

func (l *DeviceList) Layout(gtx layout.Context, th *material.Theme, source bool) layout.Dimensions {
	defer func() {
		recover()
	}()
	list := material.List(th, l.LayoutList)
	l.LayoutList.Axis = layout.Vertical
	list.AnchorStrategy = material.Overlay
	return layout.Inset{
		Top:    0,
		Bottom: 0,
		Left:   0,
		Right:  unit.Dp(12),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		var listTitleText string
		if source {
			listTitleText = "Захват"
		} else {
			listTitleText = "Вывод"
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				listTitle := material.Label(th, 15, listTitleText)
				listTitle.Font.Weight = font.SemiBold
				listTitle.Font.Typeface = th.Face
				listTitle.Alignment = text.Middle
				var listTitleWidget layout.Widget = listTitle.Layout
				return layout.UniformInset(4).Layout(gtx, listTitleWidget)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				style.DrawRect(gtx, th.ContrastFg)
				return list.Layout(gtx, len(l.Devices), func(gtx layout.Context, index int) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					if source && !l.Devices[index].Chosen && !audio.Streaming && l.Devices[index].Button.Clicked(gtx) {
						for i := 0; i < len(l.Devices); i++ {
							if i == index {
								l.Devices[index].Chosen = true
							} else {
								l.Devices[i].Chosen = false
							}
						}
						allDevicesIndex := slices.Index(l.Devices, l.Devices[index])
						if allDevicesIndex != -1 {
							devices := make([]DeviceListItem, 0)
							for _, v := range audio.AllDevices {
								if v.Type == malgo.Capture || v.Id == l.Devices[index].Id {
									continue
								}
								devices = append(devices, DeviceListItem{Id: v.Id, Name: v.Name, Button: &widget.Clickable{}})
							}
							bell.Ring("show_receiver_devices", devices)
							bell.Ring("choose_capture_device", l.Devices[index].Id)
							l.Devices[index].Chosen = true
						}

					} else if source && l.Devices[index].Chosen && !audio.Streaming && l.Devices[index].Button.Clicked(gtx) {
						l.Devices[index].Chosen = false
						bell.Ring("remove_capture_device", l.Devices[index].Id)
						devices := make([]DeviceListItem, 0)
						bell.Ring("show_receiver_devices", devices)
					} else if !source && l.Devices[index].Id.Pointer() != nil && !l.Devices[index].Chosen && !audio.Streaming && l.Devices[index].Button.Clicked(gtx) {
						l.Devices[index].Chosen = true

						bell.Ring("choose_receiver_device", l.Devices[index].Id)
					} else if !source && l.Devices[index].Chosen && !audio.Streaming && l.Devices[index].Button.Clicked(gtx) {
						l.Devices[index].Chosen = false
						bell.Ring("remove_receiver_device", l.Devices[index].Id)

					}
					return l.Devices[index].Layout(gtx, th)
				})
			}))
	})
}

func (d *DeviceListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(40))
	gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(40))

	return DrawClickableArea(gtx, 0, d.Button, func(gtx layout.Context) layout.Dimensions {
		label := material.Label(th, unit.Sp(13), d.Name)
		label.MaxLines = 1
		if d.Button.Hovered() {
			label.Color = th.ContrastFg
			pointer.Cursor.Add(pointer.CursorPointer, gtx.Ops)

			if !d.Chosen {
				style.DrawRect(gtx, th.ContrastBg)
			}
		}
		if d.Chosen {
			style.DrawRect(gtx, color.NRGBA{169, 152, 217, 0xff})
			label.Color = th.ContrastFg
		}
		return layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
			Left:   unit.Dp(8),
			Right:  unit.Dp(8),
		}.Layout(gtx, label.Layout)
		// }).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 	return layout.Flex{}.Layout(gtx, layout.Flexed(10, label.Layout),
		// 		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		// 			button := new(widget.Clickable)
		// 			return progressBarButton(gtx, button, material.ProgressBarStyle{Color: color.NRGBA{R: 255, A: 0xFF}, Progress: 1.0}.Layout)
		// 		}))
		// })
	})
}

func DrawClickableArea(gtx layout.Context, radius int, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Button.Add(gtx.Ops)
		constraints := gtx.Constraints
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints = constraints
				return w(gtx)
			}),
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				for _, c := range button.History() {
					style.DrawInk(gtx, c)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
		)
	})
}
