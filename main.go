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
	cfg := config.New()

	capEngine, err := engine.NewCaptureEngine(cfg.SampleRate, cfg.BufferSize)
	if err != nil {
		log.Fatalf("Critical: Failed to create capture engine: %v", err)
	}
	defer capEngine.Close()

	audioChan := make(chan engine.AudioBatch, 100)

	var vad processor.VAD = processor.NewSimpleThresholdVAD(cfg.Threshold)

	stats := &processor.EngineStats{StartTime: time.Now()}

	const (
		hangoverDelay = 1500 * time.Millisecond
	)
	var lastSpeechTime time.Time
	var speechBuffer []int16
	var wasRecording bool

	go func() {
		fmt.Println("\n\033[36m[SYSTEM]\033[0m Engine Online. Listening...")
		for batch := range audioChan {
			stats.TotalFrames += uint64(batch.Size)
			samples := batch.Data[:batch.Size]

			isSpeech := vad.Process(samples)
			if isSpeech {
				lastSpeechTime = time.Now()
				if !stats.IsRecording {
					stats.IsRecording = true
				}
			}

			if stats.IsRecording {
				if time.Since(lastSpeechTime) > hangoverDelay {
					stats.IsRecording = false
				}
			}

			var peak int16
			for _, s := range samples {
				abs := s
				if abs < 0 { abs = -abs }
				if abs > peak { peak = abs }
			}

			if stats.IsRecording {
				tmp := make([]int16, len(samples))
				copy(tmp, samples)
				speechBuffer = append(speechBuffer, tmp...)
			}

			processor.RenderCLI(stats, peak, isSpeech)

			if wasRecording && !stats.IsRecording {
				buf := make([]int16, len(speechBuffer))
				copy(buf, speechBuffer)
				speechBuffer = nil

				go func(samplesForStt []int16) {
					norm := processor.Normalize(samplesForStt)
					text, _ := processor.Transcribe(norm)
					fmt.Printf("\n[STT RESULT] %s\n", text)
				}(buf)
			}

			wasRecording = stats.IsRecording

			capEngine.Pool.Put(batch.Data)
		}
	}()

	if err := capEngine.Start(audioChan); err != nil {
		log.Fatalf("Failed to start hardware: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n\033[31m[EXIT]\033[0m Shutting down gracefully...")
}