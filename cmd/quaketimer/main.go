package main

import (
	"log"

	"github.com/rocheio/quake-timer/pkg/hotkey"
)

func main() {
	m, err := hotkey.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	defer m.User32.Release()

	keys := map[int16]*hotkey.Hotkey{
		1: &hotkey.Hotkey{1, hotkey.ModAlt + hotkey.ModCtrl, 'O'},
		2: &hotkey.Hotkey{2, hotkey.ModAlt + hotkey.ModShift, 'M'},
		3: &hotkey.Hotkey{3, hotkey.ModAlt + hotkey.ModCtrl, 'X'},
		4: &hotkey.Hotkey{4, hotkey.ModAlt, '1'},
	}

	for i, k := range keys {
		err := m.RegisterHotkey(i, k)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = m.SeekHotkeyLoop()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("hotkey listener stopped")
}
