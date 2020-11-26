// Package fvad provides voice activity detection by wrapping libfvad.
//
// See https://github.com/dpirch/libfvad for installation instructions.
// See example.go in this repository for sample usage.
package fvad

// #cgo pkg-config: libfvad
// #include "fvad.h"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// A Detector detects voice activity in audio frames.
type Detector struct {
	fvad *C.Fvad
}

// NewDetector creates and initializes a voice activity detector.
// The Detector must be closed when it is no longer needed.
func NewDetector() *Detector {
	return &Detector{fvad: C.fvad_new()}
}

// Close frees resources associated with d.
func (d *Detector) Close() {
	C.fvad_free(d.fvad)
	d.fvad = nil
}

// SetMode sets the sensitivity of d.
// The higher the mode, the more restrictive in reporting speech.
// Increasing the mode reduces false positives,
// at the cost of increased false negatives.
// Valid modes are 0, 1, 2, and 3.
// The default mode is 0.
func (d *Detector) SetMode(x int) error {
	errno := C.fvad_set_mode(d.fvad, C.int(x))
	if errno != 0 {
		return fmt.Errorf("invalid mode: %v", x)
	}
	return nil
}

// SetSampleRate sets the input sample rate in Hz for d.
// Valid values are 8000, 16000, 32000 and 48000. The default is 8000.
// Note that internally all processing will be done 8000 Hz;
// input data in higher sample rates will be downsampled.
func (d *Detector) SetSampleRate(x int) error {
	errno := C.fvad_set_sample_rate(d.fvad, C.int(x))
	if errno != 0 {
		return fmt.Errorf("invalid sample rate: %v", x)
	}
	return nil
}

// Reset reinitializes d, clearing all state and
// resetting mode and sample rate to their defaults.
func (d *Detector) Reset() {
	C.fvad_reset(d.fvad)
}

// Process reports whether voice has been detected in buf.
// Only frames with a length of 10, 20 or 30 ms are supported.
// For example at 8 kHz, len(buf) must be either 80, 160 or 240.
func (d *Detector) Process(buf []int16) (voice bool, err error) {
	res := C.fvad_process(d.fvad, (*C.short)(unsafe.Pointer(&buf[0])), (C.ulong)(len(buf)))
	switch res {
	case 0:
		return false, nil
	case 1:
		return true, nil
	case -1:
		return false, errors.New("Process failed")
	}
	panic("unreachable")
}
