package processor

import (
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type Recorder struct {
	file    *os.File
	encoder *wav.Encoder
	sampleRate int
}

func NewRecorder(filename string, sampleRate int) (*Recorder, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	// .wav files need a header: 16-bit, 1 channel (Mono), PCM format
	enc := wav.NewEncoder(f, sampleRate, 16, 1, 1)

	return &Recorder{
		file:       f,
		encoder:    enc,
		sampleRate: sampleRate,
	}, nil
}

func (r *Recorder) Record(samples []int16) error {
	intSamples := make([]int, len(samples))
	for i, v := range samples {
		intSamples[i] = int(v)
	}

	buf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  r.sampleRate,
		},
		Data: intSamples,
	}

	return r.encoder.Write(buf)
}

func (r *Recorder) Close() {
	r.encoder.Close()
	r.file.Close()
}