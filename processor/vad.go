package processor

type VAD interface {
    Process(samples []int16) bool
}

type SimpleThresholdVAD struct {
    Threshold int
}

func NewSimpleThresholdVAD(threshold int) *SimpleThresholdVAD {
    return &SimpleThresholdVAD{Threshold: threshold}
}

func (v *SimpleThresholdVAD) Process(samples []int16) bool {
    var peak int16
    for _, s := range samples {
        if s > peak {
            peak = s
        }
    }
    return int(peak) >= v.Threshold
}

type VADFunc func(samples []int16) bool

func (f VADFunc) Process(samples []int16) bool { return f(samples) }
