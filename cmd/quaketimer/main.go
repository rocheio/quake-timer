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

func main() {
	m, err := hotkey.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	defer m.User32.Release()

	m.AddKey("exit", hotkey.ModAlt+hotkey.ModCtrl, 'X', func() {
		m.Exit()
	})

	m.AddKey("mega-health", hotkey.ModAlt, '1', func() {
		DoAfter(time.Second*25, func() {
			FiveSecondAlert("./audio/mega-health.wav")
		})
	})

	m.AddKey("heavy-armor", hotkey.ModAlt, '2', func() {
		DoAfter(time.Second*25, func() {
			FiveSecondAlert("./audio/heavy-armor.wav")
		})
	})

	m.AddKey("quad-damage", hotkey.ModAlt, '3', func() {
		DoAfter(time.Second*115, func() {
			FiveSecondAlert("./audio/quad-damage.wav")
		})
	})

	m.AddKey("protection", hotkey.ModAlt, '4', func() {
		DoAfter(time.Second*115, func() {
			FiveSecondAlert("./audio/protection.wav")
		})
	})

	err = m.RegisterHotkeys()
	if err != nil {
		log.Fatal(err)
	}

	err = m.SeekHotkeyLoop()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("program completed")
}
