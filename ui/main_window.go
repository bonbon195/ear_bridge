package ui

import (
	"image/color"

	"gio.tools/icons"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/bonbon195/ear_bridge/audio"
	"github.com/bonbon195/ear_bridge/ui/style"
	"github.com/nuttech/bell/v2"
)

var StartStopButton = &widget.Clickable{}

var updateDevicesButton = new(widget.Clickable)

func MainWindowLayout(gtx layout.Context, th *material.Theme, sourceDeviceList *DeviceList, receiverDeviceList *DeviceList) layout.Dimensions {
	style.DrawRect(gtx, th.Bg)
	return layout.Inset{
		Top:    unit.Dp(12),
		Bottom: unit.Dp(12),
		Left:   unit.Dp(12),
		Right:  unit.Dp(0),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {

								return sourceDeviceList.Layout(gtx, th, true)
							}),
						)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return receiverDeviceList.Layout(gtx, th, false)
							}),
						)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Label(th, 13, "Захватывать можно только устройства по умолчанию")
				t.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
				iconButton := material.IconButton(th, &widget.Clickable{}, icons.ActionInfo, "")
				iconButton.Size = 14
				iconButton.Background = color.NRGBA{169, 152, 217, 0xff}
				iconButton.Inset = layout.UniformInset(2)
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

						return layout.Flex{Alignment: layout.Start}.Layout(gtx,

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Right: 4}.Layout(gtx, iconButton.Layout)
							}),
							layout.Rigid(t.Layout),
						)

					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.Y = gtx.Dp(80)
				gtx.Constraints.Min.Y = gtx.Dp(80)
				if StartStopButton.Clicked(gtx) {
					if audio.Streaming {
						bell.Ring("stop", nil)
					} else {
						bell.Ring("start", nil)
					}
					audio.Streaming = !audio.Streaming
				}
				if StartStopButton.Hovered() {
					pointer.Cursor.Add(pointer.CursorPointer, gtx.Ops)
				}

				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))
							gtx.Constraints.Min.X = gtx.Dp(unit.Dp(200))
							gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(40))
							var text string
							if audio.Streaming {
								text = "Остановить"
							} else {
								text = "Начать"

							}
							btn := material.Button(th, StartStopButton, text)
							return btn.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(40))
							iconButton := material.IconButton(th, updateDevicesButton, icons.NavigationRefresh, "")
							iconButton.Inset = layout.UniformInset(8)
							if updateDevicesButton.Clicked(gtx) {
								bell.Ring("update_devices_list", nil)
							}
							inset := layout.Inset{Left: 4}.Layout(gtx, iconButton.Layout)

							return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return inset
							})
						}),
					)
				})
			}),
		)
	})

}
