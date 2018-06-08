package main

import (
	"log"

	"github.com/rocheio/quake-timer/pkg/audio"
	"github.com/rocheio/quake-timer/pkg/hotkey"
)

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
		err := audio.PlayFile("./audio/one_bell.mp3")
		if err != nil {
			log.Fatal(err)
		}
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
