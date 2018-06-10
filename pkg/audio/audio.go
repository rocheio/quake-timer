package audio

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func PlayFiles(paths ...string) error {
	for _, p := range paths {
		err := PlayFile(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func PlayFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	s, format, err := wav.Decode(f)
	if err != nil {
		return err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))

	select {
	case <-done:
	}

	return nil
}
