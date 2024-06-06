package eventhandler

import (
	"gioui.org/io/event"
	"github.com/gen2brain/malgo"
)

// var Events chan event.Event
var Events = make(chan event.Event)

func GetEvents() <-chan event.Event {
	return Events
}

func Send(e event.Event) {
	go func() {
		Events <- e
	}()
}

type ChooseCaptureDeviceEvent struct {
	DeviceId malgo.DeviceID
}

func (e ChooseCaptureDeviceEvent) ImplementsEvent() {}
