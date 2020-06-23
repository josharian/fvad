// +build ignore

package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
	"github.com/josharian/fvad"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// Arrange for microphone input.
	portaudio.Initialize()
	defer portaudio.Terminate()

	buf := make([]int16, 960)
	stream, err := portaudio.OpenDefaultStream(1, 0, 48000, len(buf), buf)
	check(err)
	defer stream.Close()

	// Set up voice activity detector.
	vad := fvad.NewDetector()
	defer vad.Close()
	check(vad.SetMode(3))
	check(vad.SetSampleRate(48000))

	// Process incoming audio.
	check(stream.Start())
Loop:
	for {
		check(stream.Read())
		voice, err := vad.Process(buf)
		check(err)
		if voice {
			fmt.Print("voice   \r")
		} else {
			fmt.Print("no voice\r")
		}
		select {
		case <-sig:
			break Loop
		default:
		}
	}
	check(stream.Stop())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
