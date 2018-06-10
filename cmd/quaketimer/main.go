package main

import (
	"log"
	"time"

	"github.com/rocheio/quake-timer/pkg/audio"
	"github.com/rocheio/quake-timer/pkg/hotkey"
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
	name      string
	duration  time.Duration
	audioFile string
}

func (c Cooldown) Start(t time.Time) {
	remaining := c.duration - time.Now().Sub(t)
	switch {
	case remaining > 5*time.Second:
		DoAfter(remaining-5*time.Second, func() {
			log.Printf("%s in five seconds", c.name)
			FiveSecondAlert(c.audioFile)
		})
	}
}

func main() {
	m, err := hotkey.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	defer m.User32.Release()

	m.AddKey("Exit", hotkey.ModAlt+hotkey.ModCtrl, 'X', func(t time.Time) {
		m.Exit()
	})

	cd := Cooldown{"Mega Health", 30 * time.Second, "./audio/mega-health.wav"}
	m.AddKey("Mega Health", hotkey.ModAlt, '1', cd.Start)

	m.AddKey("Heavy Armor", hotkey.ModAlt, '2', func(t time.Time) {
		DoAfter(time.Second*25, func() {
			FiveSecondAlert("./audio/heavy-armor.wav")
		})
	})

	m.AddKey("Quad Damage", hotkey.ModAlt, '3', func(t time.Time) {
		DoAfter(time.Second*115, func() {
			FiveSecondAlert("./audio/quad-damage.wav")
		})
	})

	m.AddKey("Protection", hotkey.ModAlt, '4', func(t time.Time) {
		DoAfter(time.Second*115, func() {
			FiveSecondAlert("./audio/protection.wav")
		})
	})

	err = m.RegisterHotkeys()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening for key presses... press Alt+Ctrl+X to exit")
	err = m.Listen()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("program completed")
}
