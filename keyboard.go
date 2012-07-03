
package gowl

import (
	"bytes"
)

type Keyboard struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (k *Keyboard, msg []byte)
}

//// Requests
//// Events
func (k *Keyboard) HandleEvent(opcode int16, msg []byte) {
	if k.events[opcode] != nil {
		k.events[opcode](k, msg)
	}
}

type KeyboardKeymap struct {
	Format uint32
	Fd uintptr
	Size uint32
}

func (k *Keyboard) AddKeymapListener(channel chan interface{}) {
	k.listeners[0] = append(k.listeners[0], channel)
}

func keyboard_keymap(k *Keyboard, msg []byte) {
	printEvent("keymap", msg)
	var data KeyboardKeymap
	buf := bytes.NewBuffer(msg)

	format,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Format = format

	fd,err := readUintptr(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Fd = fd

	size,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Size = size

	for _,channel := range k.listeners[0] {
		go func () { channel <- data }()
	}
}

type KeyboardEnter struct {
	Serial uint32
	Surface *Surface
	Keys []interface{}
}

func (k *Keyboard) AddEnterListener(channel chan interface{}) {
	k.listeners[1] = append(k.listeners[1], channel)
}

func keyboard_enter(k *Keyboard, msg []byte) {
	printEvent("enter", msg)
	var data KeyboardEnter
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	surfaceid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surface = getObject(surfaceid).(*Surface)
	data.Surface = surface

	keys,err := readArray(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Keys = keys

	for _,channel := range k.listeners[1] {
		go func () { channel <- data }()
	}
}

type KeyboardLeave struct {
	Serial uint32
	Surface *Surface
}

func (k *Keyboard) AddLeaveListener(channel chan interface{}) {
	k.listeners[2] = append(k.listeners[2], channel)
}

func keyboard_leave(k *Keyboard, msg []byte) {
	printEvent("leave", msg)
	var data KeyboardLeave
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	surfaceid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surface = getObject(surfaceid).(*Surface)
	data.Surface = surface

	for _,channel := range k.listeners[2] {
		go func () { channel <- data }()
	}
}

type KeyboardKey struct {
	Serial uint32
	Time uint32
	Key uint32
	State uint32
}

func (k *Keyboard) AddKeyListener(channel chan interface{}) {
	k.listeners[3] = append(k.listeners[3], channel)
}

func keyboard_key(k *Keyboard, msg []byte) {
	printEvent("key", msg)
	var data KeyboardKey
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	key,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Key = key

	state,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.State = state

	for _,channel := range k.listeners[3] {
		go func () { channel <- data }()
	}
}

type KeyboardModifiers struct {
	Serial uint32
	Mods_depressed uint32
	Mods_latched uint32
	Mods_locked uint32
	Group uint32
}

func (k *Keyboard) AddModifiersListener(channel chan interface{}) {
	k.listeners[4] = append(k.listeners[4], channel)
}

func keyboard_modifiers(k *Keyboard, msg []byte) {
	printEvent("modifiers", msg)
	var data KeyboardModifiers
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	mods_depressed,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Mods_depressed = mods_depressed

	mods_latched,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Mods_latched = mods_latched

	mods_locked,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Mods_locked = mods_locked

	group,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Group = group

	for _,channel := range k.listeners[4] {
		go func () { channel <- data }()
	}
}

func NewKeyboard() (k *Keyboard) {
	k = new(Keyboard)
	k.listeners = make(map[int16][]chan interface{}, 0)

	k.events = append(k.events, keyboard_keymap)
	k.events = append(k.events, keyboard_enter)
	k.events = append(k.events, keyboard_leave)
	k.events = append(k.events, keyboard_key)
	k.events = append(k.events, keyboard_modifiers)
	return
}

func (k *Keyboard) SetId(id int32) {
	k.id = id
}

func (k *Keyboard) Id() int32 {
	return k.id
}