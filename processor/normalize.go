package processor

import "math"

// Normalize scales the audio samples so the peak absolute value
// is close to the target amplitude (30000). It returns a new
// slice and does not modify the input.
func Normalize(samples []int16) []int16 {
    if len(samples) == 0 {
        return samples
    }

    // Use int to avoid overflow when taking abs of -32768.
    var peak int
    for _, s := range samples {
        a := int(s)
        if a < 0 {
            a = -a
        }
        if a > peak {
            peak = a
        }
    }

    // If silent or already at max positive value, return a copy.
    if peak == 0 || peak == 32767 {
        out := make([]int16, len(samples))
        copy(out, samples)
        return out
    }

    const target = 30000
    scale := float64(target) / float64(peak)

    out := make([]int16, len(samples))
    for i, s := range samples {
        // Convert to int first to avoid surprising conversions
        // for the -32768 value, then apply scale as float64.
        v := float64(int(s)) * scale
        vi := int(math.Round(v))
        if vi > 32767 {
            vi = 32767
        } else if vi < -32768 {
            vi = -32768
        }
        out[i] = int16(vi)
    }
    return out
}
