package main

import (
	"fmt"
	"log"

	"github.com/gen2brain/malgo"
)

const (
	sampleRate   = 16000
	channelCount = 1
)

func main() {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		log.Fatalf("Failed to initialize malgo context: %v", err)
	}

	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	config := malgo.DefaultDeviceConfig(malgo.Capture)
	config.Capture.Format = malgo.FormatS16
	config.Capture.Channels = channelCount
	config.SampleRate = sampleRate

	onData := func(pOutputSample, pInputSample []byte, frameCount uint32) {
		// `pInputSample` contains microphone PCM data
		if pInputSample != nil {
			fmt.Printf("Captured %d bytes\n", len(pInputSample))
		}
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onData,
	}

	device, err := malgo.InitDevice(ctx.Context, config, deviceCallbacks)
	if err != nil {
		log.Fatalf("Failed to initialize capture device: %v", err)
	}

	defer device.Uninit()
	
	if err := device.Start(); err != nil {
		log.Fatalf("Failed to start capture device: %v", err)
	}

	fmt.Println("Capturing audio... Press Ctrl+C to stop.")
	select {}
}