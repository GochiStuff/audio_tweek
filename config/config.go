package config

import "flag"

type Config struct {
    SampleRate    int
    BufferSize    int
    AudioPoolSize int
    VADMode       string
    Threshold     int
}

func New() *Config {
    cfg := &Config{}

    cfg.SampleRate = 16000
    cfg.BufferSize = 2048
    cfg.AudioPoolSize = 4
    cfg.VADMode = "threshold"
    cfg.Threshold = 500

    flag.IntVar(&cfg.SampleRate, "sample-rate", cfg.SampleRate, "Audio sample rate (Hz)")
    flag.IntVar(&cfg.BufferSize, "buffer-size", cfg.BufferSize, "Buffer size in samples")
    flag.StringVar(&cfg.VADMode, "vad", cfg.VADMode, "VAD mode: threshold|ten")
    flag.IntVar(&cfg.Threshold, "threshold", cfg.Threshold, "Amplitude threshold for simple VAD")

    flag.Parse()

    return cfg
}
