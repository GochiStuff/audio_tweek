package processor

import (
	"fmt"
	"strings"
	"time"
)

type EngineStats struct {
	TotalFrames    uint64
	TotalTime      time.Duration
	StartTime      time.Time
	DroppedBatches int
	IsRecording    bool
}

type VisualizerOptions struct {
	BarWidth    int
	MaxAmplitude int
	UseColor    bool
	ShowPercent bool
}
func DefaultVisualizerOptions() VisualizerOptions {
	return VisualizerOptions{
		BarWidth:    30,
		MaxAmplitude: 15000,
		UseColor:    true,
		ShowPercent: true,
	}
}
func RenderCLI(stats *EngineStats, currentPeak int16, vadActive bool) {
	opts := DefaultVisualizerOptions()
	render(stats, opts, currentPeak, vadActive)
}
func render(stats *EngineStats, opts VisualizerOptions, currentPeak int16, vadActive bool) {
	var peak int
	if currentPeak < 0 {
		peak = -int(currentPeak)
	} else {
		peak = int(currentPeak)
	}

	if opts.MaxAmplitude <= 0 {
		opts.MaxAmplitude = 1
	}

	fill := (peak * opts.BarWidth) / opts.MaxAmplitude
	if fill < 0 {
		fill = 0
	}
	if fill > opts.BarWidth {
		fill = opts.BarWidth
	}

	bar := strings.Repeat("█", fill) + strings.Repeat("░", opts.BarWidth-fill)

	status := "[IDLE]"
	colorStart := ""
	colorEnd := ""
	if vadActive {
		status = "[VOICE]"
		if opts.UseColor {
			colorStart = "\x1b[32m"
			colorEnd = "\x1b[0m"
		}
	} else if opts.UseColor {
		colorStart = "\x1b[2m"
		colorEnd = "\x1b[0m"
	}
	uptime := time.Since(stats.StartTime).Round(time.Second)
	percent := (peak * 100) / opts.MaxAmplitude
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	if opts.ShowPercent {
		fmt.Printf("\r\033[K%s%s%s | Peak: %5d | %s %3d%% | Uptime: %s | Frames: %d",
			colorStart, status, colorEnd, peak, bar, percent, uptime, stats.TotalFrames)
	} else {
		fmt.Printf("\r\033[K%s%s%s | Peak: %5d | %s | Uptime: %s | Frames: %d",
			colorStart, status, colorEnd, peak, bar, uptime, stats.TotalFrames)
	}
}
