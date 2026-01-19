package config

import (
    "flag"
)

// Config holds runtime configuration for the application.
type Config struct {
    SampleRate     int    // audio sample rate in Hz
    BufferSize     int    // buffer size in sample frames
    AudioPoolSize  int    // not used directly, kept for future
    VADMode        string // "threshold" or "ten"
    Threshold      int    // amplitude threshold for simple VAD
    RecorderFile   string // filename for recordings
    RecorderEnable bool
}

// New returns a Config populated from flags and environment variables.
// Flags take precedence over environment variables.
func New() *Config {
    cfg := &Config{}

    // Defaults
    cfg.SampleRate = 16000
    cfg.BufferSize = 2048
    cfg.AudioPoolSize = 4
    cfg.VADMode = "threshold"
    cfg.Threshold = 500
    cfg.RecorderFile = "recording.wav"
    cfg.RecorderEnable = false

    // Flags
    flag.IntVar(&cfg.SampleRate, "sample-rate", cfg.SampleRate, "Audio sample rate (Hz)")
    flag.IntVar(&cfg.BufferSize, "buffer-size", cfg.BufferSize, "Buffer size in samples")
    flag.StringVar(&cfg.VADMode, "vad", cfg.VADMode, "VAD mode: threshold|ten")
    flag.IntVar(&cfg.Threshold, "threshold", cfg.Threshold, "Amplitude threshold for simple VAD")
    flag.StringVar(&cfg.RecorderFile, "recfile", cfg.RecorderFile, "Recorder output filename")
    flag.BoolVar(&cfg.RecorderEnable, "record", cfg.RecorderEnable, "Enable recording to WAV file")

    flag.Parse()

    return cfg
}
