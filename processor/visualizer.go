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
	BarWidth    int  // number of characters in the amplitude bar
	MaxAmplitude int // amplitude mapped to full bar
	UseColor    bool // colorize status when VAD active
	ShowPercent bool // show percentage value next to bar
}

// DefaultVisualizerOptions returns sensible defaults.
func DefaultVisualizerOptions() VisualizerOptions {
	return VisualizerOptions{
		BarWidth:    30,
		MaxAmplitude: 15000,
		UseColor:    true,
		ShowPercent: true,
	}
}

// RenderCLI keeps the original function signature for compatibility but
// uses improved rendering logic (absolute peak, colors, percent).
func RenderCLI(stats *EngineStats, currentPeak int16, vadActive bool) {
	opts := DefaultVisualizerOptions()
	render(stats, opts, currentPeak, vadActive)
}

// render draws the status line according to options.
func render(stats *EngineStats, opts VisualizerOptions, currentPeak int16, vadActive bool) {
	// work with absolute peak value
	var peak int
	if currentPeak < 0 {
		peak = -int(currentPeak)
	} else {
		peak = int(currentPeak)
	}

	if opts.MaxAmplitude <= 0 {
		opts.MaxAmplitude = 1
	}

	// compute fill (clamped)
	fill := (peak * opts.BarWidth) / opts.MaxAmplitude
	if fill < 0 {
		fill = 0
	}
	if fill > opts.BarWidth {
		fill = opts.BarWidth
	}

	bar := strings.Repeat("█", fill) + strings.Repeat("░", opts.BarWidth-fill)

	// status label
	status := "[IDLE]"
	colorStart := ""
	colorEnd := ""
	if vadActive {
		status = "[VOICE]"
		if opts.UseColor {
			// green when voice
			colorStart = "\x1b[32m"
			colorEnd = "\x1b[0m"
		}
	} else if opts.UseColor {
		// dim idle
		colorStart = "\x1b[2m"
		colorEnd = "\x1b[0m"
	}

	// uptime
	uptime := time.Since(stats.StartTime).Round(time.Second)

	// percentage (optional)
	percent := (peak * 100) / opts.MaxAmplitude
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	// print one-line dashboard (clear line first)
	if opts.ShowPercent {
		fmt.Printf("\r\033[K%s%s%s | Peak: %5d | %s %3d%% | Uptime: %s | Frames: %d",
			colorStart, status, colorEnd, peak, bar, percent, uptime, stats.TotalFrames)
	} else {
		fmt.Printf("\r\033[K%s%s%s | Peak: %5d | %s | Uptime: %s | Frames: %d",
			colorStart, status, colorEnd, peak, bar, uptime, stats.TotalFrames)
	}
}
