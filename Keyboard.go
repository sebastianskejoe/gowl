package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Keyboard struct {
	id int32
    keymapListeners []chan KeyboardKeymap
	enterListeners []chan KeyboardEnter
	leaveListeners []chan KeyboardLeave
	keyListeners []chan KeyboardKey
	modifiersListeners []chan KeyboardModifiers
	events []func(k *Keyboard, msg message) error
    name string
}

func NewKeyboard() (k *Keyboard) {
	k = new(Keyboard)
    k.name = "Keyboard"
    k.keymapListeners = make([]chan KeyboardKeymap, 0)
	k.enterListeners = make([]chan KeyboardEnter, 0)
	k.leaveListeners = make([]chan KeyboardLeave, 0)
	k.keyListeners = make([]chan KeyboardKey, 0)
	k.modifiersListeners = make([]chan KeyboardModifiers, 0)

    k.events = append(k.events, keyboardKeymap)
	k.events = append(k.events, keyboardEnter)
	k.events = append(k.events, keyboardLeave)
	k.events = append(k.events, keyboardKey)
	k.events = append(k.events, keyboardModifiers)
	return
}

func (k *Keyboard) HandleEvent(msg message) {
	if k.events[msg.opcode] != nil {
		k.events[msg.opcode](k, msg)
	}
}

func (k *Keyboard) SetId(id int32) {
	k.id = id
}

func (k *Keyboard) Id() int32 {
	return k.id
}

func (k *Keyboard) Name() string {
    return k.name
}

////
//// REQUESTS
////

////
//// EVENTS
////

type KeyboardKeymap struct {
    Format uint32
	Fd uintptr
	Size uint32
}

func (k *Keyboard) AddKeymapListener(channel chan KeyboardKeymap) {
    k.keymapListeners = append(k.keymapListeners, channel)
}

func keyboardKeymap(k *Keyboard, msg message) (err error) {
    var data KeyboardKeymap

    // Read format
    format,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Format = format

    // Read fd
    fd,err := msg.fd, nil
    if err != nil {
        return
    }
    data.Fd = fd

    // Read size
    size,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Size = size

    // Dispatch events
    for _,channel := range k.keymapListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type KeyboardEnter struct {
    Serial uint32
	Surface *Surface
	Keys []interface{}
}

func (k *Keyboard) AddEnterListener(channel chan KeyboardEnter) {
    k.enterListeners = append(k.enterListeners, channel)
}

func keyboardEnter(k *Keyboard, msg message) (err error) {
    var data KeyboardEnter

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Read surface
    surface,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    surfaceObj := getObject(surface)
    data.Surface = surfaceObj.(*Surface)

    // Read keys
    keys,err := readArray(msg.buf)
    if err != nil {
        return
    }
    data.Keys = keys

    // Dispatch events
    for _,channel := range k.enterListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type KeyboardLeave struct {
    Serial uint32
	Surface *Surface
}

func (k *Keyboard) AddLeaveListener(channel chan KeyboardLeave) {
    k.leaveListeners = append(k.leaveListeners, channel)
}

func keyboardLeave(k *Keyboard, msg message) (err error) {
    var data KeyboardLeave

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Read surface
    surface,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    surfaceObj := getObject(surface)
    data.Surface = surfaceObj.(*Surface)

    // Dispatch events
    for _,channel := range k.leaveListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type KeyboardKey struct {
    Serial uint32
	Time uint32
	Key uint32
	State uint32
}

func (k *Keyboard) AddKeyListener(channel chan KeyboardKey) {
    k.keyListeners = append(k.keyListeners, channel)
}

func keyboardKey(k *Keyboard, msg message) (err error) {
    var data KeyboardKey

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Read time
    time,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Time = time

    // Read key
    key,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Key = key

    // Read state
    state,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.State = state

    // Dispatch events
    for _,channel := range k.keyListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type KeyboardModifiers struct {
    Serial uint32
	ModsDepressed uint32
	ModsLatched uint32
	ModsLocked uint32
	Group uint32
}

func (k *Keyboard) AddModifiersListener(channel chan KeyboardModifiers) {
    k.modifiersListeners = append(k.modifiersListeners, channel)
}

func keyboardModifiers(k *Keyboard, msg message) (err error) {
    var data KeyboardModifiers

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Read mods_depressed
    mods_depressed,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.ModsDepressed = mods_depressed

    // Read mods_latched
    mods_latched,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.ModsLatched = mods_latched

    // Read mods_locked
    mods_locked,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.ModsLocked = mods_locked

    // Read group
    group,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Group = group

    // Dispatch events
    for _,channel := range k.modifiersListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
