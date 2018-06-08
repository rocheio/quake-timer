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

	// register 'exit' key as ALT + CTRL + X
	exitKey := hotkey.Hotkey{1, hotkey.ModAlt + hotkey.ModCtrl, 'X', func() {
		m.Exit()
	}}
	err = m.RegisterHotkey(1, &exitKey)
	if err != nil {
		log.Fatal(err)
	}

	// register 'mega health' key as ALT + 1
	megaHealthKey := hotkey.Hotkey{2, hotkey.ModAlt, '1', func() {
		DoAfter(time.Second*2, func() {
			FiveSecondAlert("./audio/mega-health.wav")
		})
	}}
	err = m.RegisterHotkey(2, &megaHealthKey)
	if err != nil {
		log.Fatal(err)
	}

	// loop until error or exit code, playing actions on hotkey
	err = m.SeekHotkeyLoop()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("program completed")
}
