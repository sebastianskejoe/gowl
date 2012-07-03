
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
	format uint32
	fd uintptr
	size uint32
}

func (k *Keyboard) AddKeymapListener(channel chan interface{}) {
	k.addListener(0, channel)
}

func keyboard_keymap(k *Keyboard, msg []byte) {
	printEvent("keymap", msg)
	var data KeyboardKeymap
	buf := bytes.NewBuffer(msg)

	format,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.format = format

	fd,err := readUintptr(buf)
	if err != nil {
		// XXX Error handling
	}
	data.fd = fd

	size,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.size = size

	for _,channel := range k.listeners[0] {
		channel <- data
	}
}

type KeyboardEnter struct {
	serial uint32
	surface *Surface
	keys []interface{}
}

func (k *Keyboard) AddEnterListener(channel chan interface{}) {
	k.addListener(1, channel)
}

func keyboard_enter(k *Keyboard, msg []byte) {
	printEvent("enter", msg)
	var data KeyboardEnter
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.serial = serial

	surfaceid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surface = getObject(surfaceid).(*Surface)
	data.surface = surface

	keys,err := readArray(buf)
	if err != nil {
		// XXX Error handling
	}
	data.keys = keys

	for _,channel := range k.listeners[1] {
		channel <- data
	}
}

type KeyboardLeave struct {
	serial uint32
	surface *Surface
}

func (k *Keyboard) AddLeaveListener(channel chan interface{}) {
	k.addListener(2, channel)
}

func keyboard_leave(k *Keyboard, msg []byte) {
	printEvent("leave", msg)
	var data KeyboardLeave
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.serial = serial

	surfaceid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surface = getObject(surfaceid).(*Surface)
	data.surface = surface

	for _,channel := range k.listeners[2] {
		channel <- data
	}
}

type KeyboardKey struct {
	serial uint32
	time uint32
	key uint32
	state uint32
}

func (k *Keyboard) AddKeyListener(channel chan interface{}) {
	k.addListener(3, channel)
}

func keyboard_key(k *Keyboard, msg []byte) {
	printEvent("key", msg)
	var data KeyboardKey
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.serial = serial

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.time = time

	key,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.key = key

	state,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.state = state

	for _,channel := range k.listeners[3] {
		channel <- data
	}
}

type KeyboardModifiers struct {
	serial uint32
	mods_depressed uint32
	mods_latched uint32
	mods_locked uint32
	group uint32
}

func (k *Keyboard) AddModifiersListener(channel chan interface{}) {
	k.addListener(4, channel)
}

func keyboard_modifiers(k *Keyboard, msg []byte) {
	printEvent("modifiers", msg)
	var data KeyboardModifiers
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.serial = serial

	mods_depressed,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.mods_depressed = mods_depressed

	mods_latched,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.mods_latched = mods_latched

	mods_locked,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.mods_locked = mods_locked

	group,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.group = group

	for _,channel := range k.listeners[4] {
		channel <- data
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