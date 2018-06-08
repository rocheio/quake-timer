package main

import (
	"bytes"
	"fmt"
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

type ProcManager struct {
	user32  *syscall.DLL
	regkey  *syscall.Proc
	peekmsg *syscall.Proc
}

func NewProcManager() ProcManager {
	m := ProcManager{}
	m.user32 = syscall.MustLoadDLL("user32")
	m.regkey = m.user32.MustFindProc("RegisterHotKey")
	m.peekmsg = m.user32.MustFindProc("PeekMessageW")
	return m
}

func (m ProcManager) RegisterHotkey(h *Hotkey) {
	r1, _, err := m.regkey.Call(
		0, uintptr(h.Id), uintptr(h.Modifiers), uintptr(h.KeyCode),
	)
	if r1 == 1 {
		fmt.Println("Registered", h)
	} else {
		fmt.Println("Failed to register", h, ", error:", err)
	}
}

type HotkeyMessage struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
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

func main() {
	m := NewProcManager()
	defer m.user32.Release()

	hotkeys := map[int16]*Hotkey{
		1: &Hotkey{1, ModAlt + ModCtrl, 'O'},  // ALT+CTRL+O
		2: &Hotkey{2, ModAlt + ModShift, 'M'}, // ALT+SHIFT+M
		3: &Hotkey{3, ModAlt + ModCtrl, 'X'},  // ALT+CTRL+X
		4: &Hotkey{4, ModAlt, '1'},            // F1
	}

	for _, h := range hotkeys {
		m.RegisterHotkey(h)
	}

	for {
		var msg = &HotkeyMessage{}
		m.peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			fmt.Println("Hotkey pressed:", hotkeys[id])
			if id == 3 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				return
			}
			if id == 4 { // ALT+1 = Ding
				fmt.Println("Play audio here...")
			}
		}

		time.Sleep(time.Millisecond * 50)
	}
}
