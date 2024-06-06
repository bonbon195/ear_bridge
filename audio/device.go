package audio

import (
	"log"

	"github.com/gen2brain/malgo"
)

/*
Класс для аудиоустройства
*/
type Device struct {
	Id              malgo.DeviceID
	MalgoDevice     *malgo.Device
	Name            string
	Type            malgo.DeviceType
	Config          malgo.DeviceConfig
	Callbacks       malgo.DeviceCallbacks
	CapturedSamples []byte
	ReceiverChannel chan []byte
}

// Вывод имени устройства
func (d *Device) String() string {
	return d.Name
}

// Включение устройства для обработки ивентов
func (d *Device) Start(ctx malgo.Context) error {
	d.InitMalgoDevice(ctx)
	d.MalgoDevice.Start()
	for {
		select {
		case <-Context.Done():
			d.MalgoDevice.Uninit()
			return nil
		}
	}

}

// Инициализация устройства
func (d *Device) InitMalgoDevice(ctx malgo.Context) error {
	var err error
	d.MalgoDevice, err = malgo.InitDevice(ctx, d.Config, d.Callbacks)
	if err != nil {
		return err
	}
	return nil
}

// Инициализация обратных вызовов для устройства захвата.
// Копирует данные из буфера в каналы устройств вывода.
// Для каждого копирования запускается горутина
func (d *Device) InitCaptureCallbacks() {
	d.Callbacks = malgo.DeviceCallbacks{
		Data: func(pOutputSample, pInputSamples []byte, framecount uint32) {
			if framecount == 0 {
				return
			}
			// for _, v := range CurrentReceiverDevices {
			// 	go func(v *Device) { v.ReceiverChannel <- pInputSamples }(v)
			// }
			d.CapturedSamples = pInputSamples
		},
		Stop: func() {},
	}
}

// Инициалиация обратных вызовов для устройства вывода
// Копирует данные из канала в буфер.
// Поток не блокируется, даже если данных в канале нет.
func (d *Device) InitReceiverCallbacks(sourceDevice *Device) {
	d.ReceiverChannel = make(chan []byte)
	d.Callbacks = malgo.DeviceCallbacks{
		Data: func(pOutputSample, pInputSample []byte, framecount uint32) {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()

			if sourceDevice.CapturedSamples != nil && len(sourceDevice.CapturedSamples) > 0 {
				copy(pOutputSample, sourceDevice.CapturedSamples)
			}
			// select {
			// case data := <-d.ReceiverChannel:
			// 	copy(pOutputSample, data)
			// default:
			// }

		},
	}
}

// Установка начальной конфигурации устройства
// Следует вызывать перед вызовом метода Start()
func (d *Device) SetDeviceConfig(dType malgo.DeviceType) {
	d.Config = malgo.DefaultDeviceConfig(dType)
	setConfigOpts(&d.Config)
	d.Config.Playback.DeviceID = d.Id.Pointer()
}

func setConfigOpts(config *malgo.DeviceConfig) {
	config.Alsa.NoMMap = 1
	config.NoClip = 1
	config.PerformanceProfile = malgo.LowLatency
	config.Capture.Channels = 2
	config.Playback.Channels = 2
	config.SampleRate = 44100
	config.Capture.Format = malgo.FormatS16
	config.Playback.Format = malgo.FormatS16
	// config.SampleRate = 192000
	// config.Capture.Format = malgo.FormatS24
	// config.Playback.Format = malgo.FormatS24
	// config.Resampling.Algorithm = malgo.ResampleAlgorithmSpeex
}
