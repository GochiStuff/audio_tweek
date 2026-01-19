package engine

import (
	"sync"
	"unsafe"

	"github.com/gen2brain/malgo"
)

type AudioBatch struct {
	Data []int16
	Size int
}

type CaptureEngine struct {
	ctx    *malgo.AllocatedContext
	device *malgo.Device
	Pool   *sync.Pool
	SampleRate int
	BufferSize int
}

func NewCaptureEngine(sampleRate, bufferSize int) (*CaptureEngine, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}

	if bufferSize <= 0 {
		bufferSize = 2048
	}

	pool := &sync.Pool{
		New: func() interface{} {
			return make([]int16, bufferSize)
		},
	}

	return &CaptureEngine{ctx: ctx, Pool: pool, SampleRate: sampleRate, BufferSize: bufferSize}, nil
}

func (e *CaptureEngine) Start(audioChan chan<- AudioBatch) error {
	config := malgo.DefaultDeviceConfig(malgo.Capture)
	config.Capture.Format = malgo.FormatS16
	config.Capture.Channels = 1
	if e.SampleRate > 0 {
		config.SampleRate = uint32(e.SampleRate)
	}

	onData := func(pOutput, pInput []byte, frameCount uint32) {
		if pInput == nil {
			return
		}

		sampleCount := len(pInput) / 2
		rawSamples := unsafe.Slice((*int16)(unsafe.Pointer(&pInput[0])), sampleCount)

		buf := e.Pool.Get().([]int16)
		if len(buf) < sampleCount {
			buf = make([]int16, sampleCount)
		}
		copy(buf, rawSamples)

		audioChan <- AudioBatch{Data: buf, Size: sampleCount}
	}

	device, err := malgo.InitDevice(e.ctx.Context, config, malgo.DeviceCallbacks{Data: onData})
	if err != nil {
		return err
	}
	e.device = device

	return e.device.Start()
}

func (e *CaptureEngine) Close() {
	if e.device != nil {
		e.device.Uninit()
	}
	if e.ctx != nil {
		_ = e.ctx.Uninit()
		e.ctx.Free()
	}
}