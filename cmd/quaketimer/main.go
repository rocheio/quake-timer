package main

import (
	"log"
	"time"

	"github.com/rocheio/quake-timer/pkg/cooldown"
	"github.com/rocheio/quake-timer/pkg/hotkey"
)

func main() {
	m, err := hotkey.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	defer m.User32.Release()

	m.AddKey("Exit", hotkey.ModAlt+hotkey.ModCtrl, 'X', func(t time.Time) {
		m.Exit()
	})

	cd := cooldown.Cooldown{"Mega Health", 30 * time.Second, "./audio/mega-health.wav"}
	m.AddKey("Mega Health", hotkey.ModAlt, '1', cd.Start)

	cd = cooldown.Cooldown{"Heavy Armor", 30 * time.Second, "./audio/heavy-armor.wav"}
	m.AddKey("Heavy Armor", hotkey.ModAlt, '2', cd.Start)

	cd = cooldown.Cooldown{"Quad Damage", 120 * time.Second, "./audio/quad-damage.wav"}
	m.AddKey("Quad Damage", hotkey.ModAlt, '3', cd.Start)

	cd = cooldown.Cooldown{"Protection", 120 * time.Second, "./audio/protection.wav"}
	m.AddKey("Protection", hotkey.ModAlt, '4', cd.Start)

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
