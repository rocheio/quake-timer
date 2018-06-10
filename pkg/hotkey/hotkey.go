package hotkey

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

type Manager struct {
	User32  *syscall.DLL
	regkey  *syscall.Proc
	peekmsg *syscall.Proc
	keys    map[int16]*Hotkey
	exit    bool
}

func NewManager() (*Manager, error) {
	var err error
	m := Manager{
		keys: make(map[int16]*Hotkey),
	}
	m.User32, err = syscall.LoadDLL("user32")
	if err != nil {
		return nil, err
	}
	m.regkey, err = m.User32.FindProc("RegisterHotKey")
	if err != nil {
		return nil, err
	}
	m.peekmsg, err = m.User32.FindProc("PeekMessageW")
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Manager) RegisterHotkey(i int16, h *Hotkey) error {
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

func (m *Manager) RegisterHotkeys(keys map[string]*Hotkey) error {
	for _, key := range keys {
		err := m.RegisterHotkey(int16(key.Id), key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) SeekHotkeyID() (int16, error) {
	msg := &Message{}
	m.peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)
	return msg.WPARAM, nil
}

func (m *Manager) SeekHotkeyLoop() error {
	for ; ; time.Sleep(time.Millisecond * 50) {
		if m.exit {
			log.Println("received signal to exit")
			return nil
		}

		id, err := m.SeekHotkeyID()
		if err != nil {
			return err
		}
		if id == 0 {
			continue
		}

		key := m.keys[id]
		log.Println("Hotkey pressed:", key)

		if key.Action != nil {
			go key.Action()
		}
	}
}

func (m *Manager) Exit() {
	m.exit = true
}

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
	Action    func()
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

type Message struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}
