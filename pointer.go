package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Pointer struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (p *Pointer, msg message)
}

//// Requests
func (p *Pointer) SetCursor (serial uint32, surface *Surface, hotspot_x int32, hotspot_y int32) {
	msg := newMessage(p, 0)
	writeInteger(msg,serial)
	writeInteger(msg,surface.Id())
	writeInteger(msg,hotspot_x)
	writeInteger(msg,hotspot_y)

	sendmsg(msg)
	printRequest("pointer", p, "set_cursor", serial, surface.Id(), hotspot_x, hotspot_y)
}

//// Events
func (p *Pointer) HandleEvent(msg message) {
	if p.events[msg.opcode] != nil {
		p.events[msg.opcode](p, msg)
	}
}

type PointerEnter struct {
	Serial uint32
	Surface *Surface
	SurfaceX int32
	SurfaceY int32
}

func (p *Pointer) AddEnterListener(channel chan interface{}) {
	p.listeners[0] = append(p.listeners[0], channel)
	addListener(channel)
}

func pointer_enter(p *Pointer, msg message) {
	var data PointerEnter

	serial,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	surfaceid, err := readInt32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surfaceobj := getObject(surfaceid)
	if surfaceobj == nil {
		return
	}
	surface = surfaceobj.(*Surface)
	data.Surface = surface

	surface_x,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.SurfaceX = surface_x

	surface_y,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.SurfaceY = surface_y

	for _,channel := range p.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("pointer", p, "enter", serial, surface.Id(), surface_x, surface_y)
}

type PointerLeave struct {
	Serial uint32
	Surface *Surface
}

func (p *Pointer) AddLeaveListener(channel chan interface{}) {
	p.listeners[1] = append(p.listeners[1], channel)
	addListener(channel)
}

func pointer_leave(p *Pointer, msg message) {
	var data PointerLeave

	serial,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	surfaceid, err := readInt32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surfaceobj := getObject(surfaceid)
	if surfaceobj == nil {
		return
	}
	surface = surfaceobj.(*Surface)
	data.Surface = surface

	for _,channel := range p.listeners[1] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("pointer", p, "leave", serial, surface.Id())
}

type PointerMotion struct {
	Time uint32
	SurfaceX int32
	SurfaceY int32
}

func (p *Pointer) AddMotionListener(channel chan interface{}) {
	p.listeners[2] = append(p.listeners[2], channel)
	addListener(channel)
}

func pointer_motion(p *Pointer, msg message) {
	var data PointerMotion

	time,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	surface_x,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.SurfaceX = surface_x

	surface_y,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.SurfaceY = surface_y

	for _,channel := range p.listeners[2] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("pointer", p, "motion", time, surface_x, surface_y)
}

type PointerButton struct {
	Serial uint32
	Time uint32
	Button uint32
	State uint32
}

func (p *Pointer) AddButtonListener(channel chan interface{}) {
	p.listeners[3] = append(p.listeners[3], channel)
	addListener(channel)
}

func pointer_button(p *Pointer, msg message) {
	var data PointerButton

	serial,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	time,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	button,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Button = button

	state,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.State = state

	for _,channel := range p.listeners[3] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("pointer", p, "button", serial, time, button, state)
}

type PointerAxis struct {
	Time uint32
	Axis uint32
	Value int32
}

func (p *Pointer) AddAxisListener(channel chan interface{}) {
	p.listeners[4] = append(p.listeners[4], channel)
	addListener(channel)
}

func pointer_axis(p *Pointer, msg message) {
	var data PointerAxis

	time,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	axis,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Axis = axis

	value,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Value = value

	for _,channel := range p.listeners[4] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("pointer", p, "axis", time, axis, value)
}

func NewPointer() (p *Pointer) {
	p = new(Pointer)
	p.listeners = make(map[int16][]chan interface{}, 0)

	p.events = append(p.events, pointer_enter)
	p.events = append(p.events, pointer_leave)
	p.events = append(p.events, pointer_motion)
	p.events = append(p.events, pointer_button)
	p.events = append(p.events, pointer_axis)
	return
}

func (p *Pointer) SetId(id int32) {
	p.id = id
}

func (p *Pointer) Id() int32 {
	return p.id
}