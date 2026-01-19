package processor

// NormalizeInPlace scales samples in-place by delegating to Normalize
// and copying the result back into the provided slice.
func NormalizeInPlace(samples []int16) []int16 {
    out := Normalize(samples)
    copy(samples, out)
    return samples
}