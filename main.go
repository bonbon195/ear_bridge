package main

import (
	"context"
	"image/color"
	"log"
	"os"
	"sync"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/bonbon195/ear_bridge/audio"
	"github.com/bonbon195/ear_bridge/ui"
	"github.com/gen2brain/malgo"
	"github.com/nuttech/bell/v2"
	"golang.org/x/exp/slices"
)

var sourceDeviceList = ui.DeviceList{LayoutList: &widget.List{}}
var receiverDeviceList = ui.DeviceList{LayoutList: &widget.List{}}

func main() {

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()
	go func() {
		w := app.NewWindow(app.Title("EarBridge"))
		err := run(w, ctx)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	updateDevicesList(ctx)

	app.Main()
}

func run(w *app.Window, ctx *malgo.AllocatedContext) error {
	th := material.NewTheme()
	th.ContrastBg = color.NRGBA{138, 118, 192, 0xff}
	th.Bg = color.NRGBA{240, 238, 247, 0xff}
	var ops op.Ops
	bell.Listen("update_devices_list", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		updateDevicesList(ctx)
	})
	bell.Listen("show_receiver_devices", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		receiverDeviceList.Devices = message.([]ui.DeviceListItem)
	})

	bell.Listen("choose_capture_device", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		id := message.(malgo.DeviceID)
		for i := 0; i < len(audio.AllDevices); i++ {
			if audio.AllDevices[i].Id == id {

				audio.CurrentCaptureDevice = &audio.AllDevices[i]

				for i, v := range audio.CurrentReceiverDevices {
					if v.Id == audio.CurrentCaptureDevice.Id {
						audio.CurrentReceiverDevices = slices.Delete(audio.CurrentReceiverDevices, i, i+1)
					}
				}

				var dType malgo.DeviceType
				if audio.CurrentCaptureDevice.Type == malgo.Playback {
					dType = malgo.Loopback
				} else {
					dType = malgo.Capture
				}
				audio.CurrentCaptureDevice.SetDeviceConfig(dType)
				// audio.CurrentCaptureDevice.Device.Config.Capture.Format = malgo.FormatS24
				// audio.CurrentCaptureDevice.Device.Config.SampleRate = 192000
				audio.CurrentCaptureDevice.InitCaptureCallbacks()
				break
			}
		}
	})

	bell.Listen("remove_capture_device", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		audio.CurrentCaptureDevice = nil

	})

	bell.Listen("choose_receiver_device", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		id := message.(malgo.DeviceID)

		for i := 0; i < len(audio.CurrentReceiverDevices); i++ {
			if audio.CurrentReceiverDevices[i].Id == id {
				return
			}
		}
		for i := 0; i < len(audio.AllDevices); i++ {
			if audio.AllDevices[i].Id == id {
				if audio.AllDevices[i].Type == malgo.Capture {
					audio.AllDevices[i].SetDeviceConfig(malgo.Capture)
				} else {
					audio.AllDevices[i].SetDeviceConfig(malgo.Playback)
					// audio.AllDevices[i].Config.SampleRate = 192000
					// audio.AllDevices[i].Config.Playback.Format = malgo.FormatS24
				}
				audio.AllDevices[i].InitReceiverCallbacks(audio.CurrentCaptureDevice)
				audio.CurrentReceiverDevices = append(audio.CurrentReceiverDevices, &audio.AllDevices[i])
				break
			}
		}
	})

	bell.Listen("remove_receiver_device", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		id := message.(malgo.DeviceID)

		for i, v := range audio.CurrentReceiverDevices {
			if v.Id == id {
				if len(audio.CurrentReceiverDevices) == 1 {
					audio.Streaming = false
				}
				audio.CurrentReceiverDevices = slices.Delete(audio.CurrentReceiverDevices, i, i+1)
				break
			}
		}

	})

	bell.Listen("start", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		audio.Context, audio.Cancel = context.WithCancel(context.Background())
		if audio.CurrentCaptureDevice.Id.Pointer() != nil && len(audio.CurrentReceiverDevices) != 0 {
			go func() {
				audio.CurrentCaptureDevice.Start(ctx.Context)
			}()

			for i := 0; i < len(audio.CurrentReceiverDevices); i++ {
				go func(i int) {
					audio.CurrentReceiverDevices[i].Start(ctx.Context)
				}(i)

			}

		}
	})

	bell.Listen("stop", func(message bell.Message) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		audio.Cancel()
		for _, v := range audio.CurrentReceiverDevices {
			go func(v *audio.Device) {
				// TODO
			}(v)
		}
	})
	// for {
	// 	select {
	// 	case e := <-w.Events():
	// 		switch e := e.(type) {
	// 		case system.DestroyEvent:
	// 			return e.Err
	// 		case system.FrameEvent:
	// 			gtx := layout.NewContext(&ops, e)
	// 			defer func() {
	// 				if r := recover(); r != nil {
	// 					log.Println(r)
	// 				}
	// 			}()
	// 			ui.MainWindowLayout(gtx, th, &sourceDeviceList, &receiverDeviceList)
	// 			e.Frame(gtx.Ops)
	// 		}
	// 		// Handle system and user interaction events
	// 	// case e := <-eventhandler.Events:
	// 	// 	switch e := e.(type) {
	// 	// 	case eventhandler.ChooseCaptureDeviceEvent:
	// 	// 		onChooseCaptureDevice(e.DeviceId)
	// 	// 		// Handle custom events
	// 	// 	}
	// 	// }
	// }
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.MainWindowLayout(gtx, th, &sourceDeviceList, &receiverDeviceList)
			e.Frame(gtx.Ops)
		}
	}
}

func onChooseCaptureDevice(id malgo.DeviceID) {
	for i := 0; i < len(audio.CurrentReceiverDevices); i++ {
		if audio.CurrentReceiverDevices[i].Id == id {
			return
		}
	}
	for i := 0; i < len(audio.AllDevices); i++ {
		if audio.AllDevices[i].Id == id {
			if audio.AllDevices[i].Type == malgo.Capture {
				audio.AllDevices[i].SetDeviceConfig(malgo.Capture)
			} else {
				audio.AllDevices[i].SetDeviceConfig(malgo.Playback)
				// audio.AllDevices[i].Config.SampleRate = 192000
				// audio.AllDevices[i].Config.Playback.Format = malgo.FormatS24

			}
			audio.AllDevices[i].InitReceiverCallbacks(audio.CurrentCaptureDevice)
			audio.CurrentReceiverDevices = append(audio.CurrentReceiverDevices, &audio.AllDevices[i])
			break
		}
	}
}

func updateDevicesList(ctx *malgo.AllocatedContext) {
	outputDevices, err := ctx.Devices(malgo.Playback)
	if err != nil {
		log.Fatal(err)
	}
	captureDevices, err := ctx.Devices(malgo.Capture)
	if err != nil {
		log.Fatal(err)
	}
	mutex := &sync.Mutex{}
	w := &sync.WaitGroup{}
	allDevices := make([]audio.Device, 0)
	sourceDevices := make([]ui.DeviceListItem, 0)
	receiverDeviceList.Devices = nil
	w.Add(2)
	go func(mutex *sync.Mutex) {
		mutex.Lock()
		for _, v := range outputDevices {

			device := audio.Device{Id: v.ID, Name: v.Name(), Type: malgo.Playback}
			allDevices = append(allDevices, device)
			// if v.IsDefault == 1 {
			sourceDevices = append(sourceDevices, ui.DeviceListItem{Id: v.ID, Name: v.Name(), Button: &widget.Clickable{}})
			// }
		}
		mutex.Unlock()
		w.Done()
	}(mutex)

	go func(mutex *sync.Mutex) {
		mutex.Lock()
		for _, v := range captureDevices {
			device := audio.Device{Id: v.ID, Name: v.Name(), Type: malgo.Capture}
			allDevices = append(allDevices, device)
			sourceDevices = append(sourceDevices, ui.DeviceListItem{Id: v.ID, Name: v.Name(), Button: &widget.Clickable{}})
		}
		mutex.Unlock()
		w.Done()
	}(mutex)
	w.Wait()
	sourceDeviceList.Devices = sourceDevices
	audio.AllDevices = allDevices
}
