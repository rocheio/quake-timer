package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

type HotkeyManager struct {
	user32  *syscall.DLL
	regkey  *syscall.Proc
	peekmsg *syscall.Proc
	keys    map[int16]*Hotkey
}

func NewHotkeyManager() (*HotkeyManager, error) {
	var err error
	m := HotkeyManager{
		keys: make(map[int16]*Hotkey),
	}
	m.user32, err = syscall.LoadDLL("user32")
	if err != nil {
		return nil, err
	}
	m.regkey, err = m.user32.FindProc("RegisterHotKey")
	if err != nil {
		return nil, err
	}
	m.peekmsg, err = m.user32.FindProc("PeekMessageW")
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m HotkeyManager) RegisterHotkey(i int16, h *Hotkey) error {
	r1, _, err := m.regkey.Call(
		0, uintptr(h.Id), uintptr(h.Modifiers), uintptr(h.KeyCode),
	)

	if r1 == 1 {
		log.Printf("registered hotkey %s", h)
		m.keys[i] = h
		return nil
	}

	return fmt.Errorf("Failed to register %v, error: %v", h, err)
}

func (m HotkeyManager) SeekHotkeyID() (int16, error) {
	msg := &HotkeyMessage{}
	m.peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)
	return msg.WPARAM, nil
}

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

type HotkeyMessage struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

func main() {
	m, err := NewHotkeyManager()
	if err != nil {
		log.Fatal(err)
	}
	defer m.user32.Release()

	keys := map[int16]*Hotkey{
		1: &Hotkey{1, ModAlt + ModCtrl, 'O'},
		2: &Hotkey{2, ModAlt + ModShift, 'M'},
		3: &Hotkey{3, ModAlt + ModCtrl, 'X'},
		4: &Hotkey{4, ModAlt, '1'},
	}

	for i, k := range keys {
		err := m.RegisterHotkey(i, k)
		if err != nil {
			log.Fatal(err)
		}
	}

loop:
	for {
		id, err := m.SeekHotkeyID()
		if err != nil {
			log.Fatal(err)
			break
		}

		switch id {
		case 0:
			break
		case 3:
			log.Println("CTRL+ALT+X pressed, goodbye...")
			break loop
		case 4:
			log.Println("ALT+1 - Play audio here...")
		default:
			log.Println("Hotkey pressed:", m.keys[id])
		}

		time.Sleep(time.Millisecond * 50)
	}
}
