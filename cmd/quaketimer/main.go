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

	hotkeys := make(map[string]*hotkey.Hotkey)
	// register 'exit' key as ALT + CTRL + X
	hotkeys["exit"] = &hotkey.Hotkey{1, hotkey.ModAlt + hotkey.ModCtrl, 'X', func() {
		m.Exit()
	}}
	// register 'mega health' key as ALT + 1
	hotkeys["mega-health"] = &hotkey.Hotkey{2, hotkey.ModAlt, '1', func() {
		DoAfter(time.Second*25, func() {
			FiveSecondAlert("./audio/mega-health.wav")
		})
	}}
	// register 'heavy armor' key as ALT + 2
	hotkeys["heavy-armor"] = &hotkey.Hotkey{3, hotkey.ModAlt, '2', func() {
		DoAfter(time.Second*25, func() {
			FiveSecondAlert("./audio/heavy-armor.wav")
		})
	}}
	// register 'quad damage' key as ALT + 3
	hotkeys["quad-damage"] = &hotkey.Hotkey{4, hotkey.ModAlt, '3', func() {
		DoAfter(time.Second*115, func() {
			FiveSecondAlert("./audio/quad-damage.wav")
		})
	}}
	// register 'protection' key as ALT + 4
	hotkeys["protection"] = &hotkey.Hotkey{5, hotkey.ModAlt, '4', func() {
		DoAfter(time.Second*115, func() {
			FiveSecondAlert("./audio/protection.wav")
		})
	}}

	err = m.RegisterHotkeys(hotkeys)
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
