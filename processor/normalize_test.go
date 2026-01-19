package processor

import "testing"

func maxAbs(samples []int16) int {
    m := 0
    for _, v := range samples {
        a := int(v)
        if a < 0 {
            a = -a
        }
        if a > m {
            m = a
        }
    }
    return m
}

func TestNormalizeSimple(t *testing.T) {
    in := []int16{100, -50, 0}
    out := Normalize(in)
    if got := maxAbs(out); got != 30000 {
        t.Fatalf("expected peak 30000, got %d", got)
    }
}

func TestNormalizeEdgeMinInt16(t *testing.T) {
    in := []int16{-32768, 0}
    out := Normalize(in)
    if len(out) != len(in) {
        t.Fatalf("unexpected length: got %d, want %d", len(out), len(in))
    }
    // Ensure values remain in int16 range and peak is scaled to target
    for i, v := range out {
        if v > 32767 || v < -32768 {
            t.Fatalf("value out of int16 range at %d: %d", i, v)
        }
    }
    if got := maxAbs(out); got != 30000 {
        t.Fatalf("expected peak 30000, got %d", got)
    }
}
