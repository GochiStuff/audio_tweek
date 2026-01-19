package processor

import "fmt"

func Transcribe(samples []int16) (string, error) {
    fmt.Printf("\n[STT] received %d samples -> processing...\n", len(samples))
    return "stt working", nil
}
