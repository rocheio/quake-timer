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
	// keys a map of hotkeys that are currently registered to the system
	keys map[int16]*Hotkey
	// exit is a flag that can be set from Goroutines to stop Listen
	exit bool
	// two attrs to track the first known Windows time registered
	// and the conversion of that int32 to the Unix time
	firstDur  int32
	firstTime time.Time
	// optional action to take each time a key press is registered
	OnKeyPress func()
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

func (m *Manager) AddKey(name string, modifiers int, keycode int, action func(d time.Time)) {
	i := len(m.keys) + 1
	m.keys[int16(i)] = &Hotkey{
		Id:        i,
		Modifiers: modifiers,
		KeyCode:   keycode,
		Name:      name,
		Action:    action,
	}
}

func (m *Manager) RegisterHotkey(i int16, h *Hotkey) error {
	r1, _, err := m.regkey.Call(
		0, uintptr(h.Id), uintptr(h.Modifiers), uintptr(h.KeyCode),
	)

	if r1 == 1 {
		log.Printf("registered %s", h)
		m.keys[i] = h
		return nil
	}

	if err.Error() == "Hot key is already registered." {
		return nil
	}

	return fmt.Errorf("Failed to register %v, error: %v", h, err)
}

func (m *Manager) RegisterHotkeys() error {
	for _, key := range m.keys {
		err := m.RegisterHotkey(int16(key.Id), key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) SeekKeyPress() (*KeyPress, error) {
	msg := &WindowsMessage{}
	m.peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

	if m.firstDur == 0 {
		m.firstDur = msg.DWORD
		m.firstTime = time.Now()
	}

	if msg.WPARAM == 0 {
		return nil, nil
	}

	key, ok := m.keys[msg.WPARAM]
	if !ok {
		return nil, fmt.Errorf("key not registered to manager: %d", msg.WPARAM)
	}

	pressTime := m.firstTime.Add(msg.DurationSince(m.firstDur))
	return &KeyPress{key, pressTime}, nil
}

func (m *Manager) Listen() error {
	for {
		select {
		case <-time.After(time.Millisecond * 100):
			if m.exit {
				log.Println("received signal to exit")
				return nil
			}

			p, err := m.SeekKeyPress()
			if err != nil {
				return err
			}
			if p == nil {
				continue
			}

			if m.OnKeyPress != nil {
				m.OnKeyPress()
			}

			log.Printf("%s pressed at %s", p.key, p.time.Format("15:04:05"))
			if p.key.Action != nil {
				go p.key.Action(p.time)
			}
		}
	}
}

func (m *Manager) Exit() {
	m.exit = true
}

type Hotkey struct {
	Id        int    // Unique id
	Modifiers int    // Mask of modifiers
	KeyCode   int    // Key code, e.g. 'A'
	Name      string // User-defined name for hotkey
	// Action does something based on the time of a keypress
	Action func(t time.Time)
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
	return fmt.Sprintf("Hotkey[%s, %s%c]", h.Name, mod, h.KeyCode)
}

// WindowsMessage is a message from a Windows thread's queue
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644958(v=vs.85).aspx
type WindowsMessage struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

func (w WindowsMessage) DurationSince(firstDur int32) time.Duration {
	return time.Millisecond * time.Duration(w.DWORD-firstDur)
}

// KeyPress is a validated press of a Hotkey at a given time
type KeyPress struct {
	key  *Hotkey
	time time.Time
}
