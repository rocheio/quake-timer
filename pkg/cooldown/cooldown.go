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

type Cooldown struct {
	Name      string
	Duration  time.Duration
	AudioFile string
}

func (c Cooldown) Start(t time.Time) {
	remaining := c.Duration - time.Now().Sub(t)
	if remaining > 20*time.Second {
		go DoAfter(remaining-20*time.Second, func() {
			log.Printf("%s in twenty seconds", c.Name)
			err := audio.PlayFiles(c.AudioFile, "./audio/in-twenty-seconds.wav")
			if err != nil {
				log.Fatal(err)
			}
		})
	}
	if remaining > 10*time.Second {
		go DoAfter(remaining-10*time.Second, func() {
			log.Printf("%s in ten seconds", c.Name)
			err := audio.PlayFiles(c.AudioFile, "./audio/in-ten-seconds.wav")
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}
