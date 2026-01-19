package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GochiStuff/audio_tweek/config"
	"github.com/GochiStuff/audio_tweek/engine"
	"github.com/GochiStuff/audio_tweek/processor"
)

func main() {
	// Load configuration from flags and environment
	cfg := config.New()

	capEngine, err := engine.NewCaptureEngine(cfg.SampleRate, cfg.BufferSize)
	if err != nil {
		log.Fatalf("Failed to create capture engine: %v", err)
	}
	defer capEngine.Close()

	audioChan := make(chan engine.AudioBatch, 100)

	if err := capEngine.Start(audioChan); err != nil {
		log.Fatalf("Failed to start capture engine: %v", err)
	}

	// recorder/stats
	stats := &processor.EngineStats{StartTime: time.Now()}

	// Choose VAD implementation based on configuration
	var vad processor.VAD
	switch cfg.VADMode {
	case "ten":
		ten, _ := processor.NewTenVADHandler()
		// NewTenVADHandler returns a handler with IsSpeech method; adapt to VAD interface
		// Wrap it so we conform to processor.VAD
		vad = processor.VADFunc(func(samples []int16) bool { return ten.IsSpeech(samples) })
	default:
		vad = processor.NewSimpleThresholdVAD(cfg.Threshold)
	}

	// Optional recorder
	var recorder *processor.Recorder
	if cfg.RecorderEnable {
		r, err := processor.NewRecorder(cfg.RecorderFile, cfg.SampleRate)
		if err != nil {
			log.Printf("Recorder init failed: %v (recording disabled)", err)
		} else {
			recorder = r
			defer recorder.Close()
			stats.IsRecording = true
		}
	}

	// worker
	go func() {
		for batch := range audioChan {
			stats.TotalFrames += uint64(batch.Size)
			samples := batch.Data[:batch.Size]
			isSpeech := vad.Process(samples)

			var peak int16
			for _, s := range samples {
				if s > peak {
					peak = s
				}
			}

			if recorder != nil && isSpeech {
				// best-effort record; ignore errors
				_ = recorder.Record(samples)
			}

			processor.RenderCLI(stats, peak, isSpeech)
			capEngine.Pool.Put(batch.Data)
		}
	}()

	fmt.Println("Capturing audio... Press Ctrl+C to stop.")

	// Keep the main thread alive until interrupted
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Stopping audio capture...")
}