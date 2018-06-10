package cooldown

import (
	"log"
	"time"

	"github.com/rocheio/quake-timer/pkg/audio"
)

func DoAfter(d time.Duration, f func()) {
	select {
	case <-time.After(d):
		f()
	}
}

func FiveSecondAlert(f string) {
	err := audio.PlayFile(f)
	if err != nil {
		log.Fatal(err)
	}
	err = audio.PlayFile("./audio/in-five-seconds.wav")
	if err != nil {
		log.Fatal(err)
	}
}

type Cooldown struct {
	Name      string
	Duration  time.Duration
	AudioFile string
}

func (c Cooldown) Start(t time.Time) {
	remaining := c.Duration - time.Now().Sub(t)
	switch {
	case remaining > 5*time.Second:
		DoAfter(remaining-5*time.Second, func() {
			log.Printf("%s in five seconds", c.Name)
			FiveSecondAlert(c.AudioFile)
		})
	}
}
