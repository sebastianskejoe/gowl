
package gowl

import (
	"bytes"
)

type Pointer struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (p *Pointer, msg []byte)
}

//// Requests
func (p *Pointer) Set_cursor (serial uint32, surface *Surface, hotspot_x int32, hotspot_y int32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, serial)
	writeInteger(buf, surface.Id())
	writeInteger(buf, hotspot_x)
	writeInteger(buf, hotspot_y)

	sendmsg(p, 0, buf.Bytes())
}

//// Events
func (p *Pointer) HandleEvent(opcode int16, msg []byte) {
	if p.events[opcode] != nil {
		p.events[opcode](p, msg)
	}
}

type PointerEnter struct {
	serial uint32
	surface *Surface
	surface_x int32
	surface_y int32
}

func (p *Pointer) AddEnterListener(channel chan interface{}) {
	p.addListener(0, channel)
}

func pointer_enter(p *Pointer, msg []byte) {
	printEvent("enter", msg)
	var data PointerEnter
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

	surface_x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.surface_x = surface_x

	surface_y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.surface_y = surface_y

	for _,channel := range p.listeners[0] {
		channel <- data
	}
}

type PointerLeave struct {
	serial uint32
	surface *Surface
}

func (p *Pointer) AddLeaveListener(channel chan interface{}) {
	p.addListener(1, channel)
}

func pointer_leave(p *Pointer, msg []byte) {
	printEvent("leave", msg)
	var data PointerLeave
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

	for _,channel := range p.listeners[1] {
		channel <- data
	}
}

type PointerMotion struct {
	time uint32
	surface_x int32
	surface_y int32
}

func (p *Pointer) AddMotionListener(channel chan interface{}) {
	p.addListener(2, channel)
}

func pointer_motion(p *Pointer, msg []byte) {
	printEvent("motion", msg)
	var data PointerMotion
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.time = time

	surface_x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.surface_x = surface_x

	surface_y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.surface_y = surface_y

	for _,channel := range p.listeners[2] {
		channel <- data
	}
}

type PointerButton struct {
	serial uint32
	time uint32
	button uint32
	state uint32
}

func (p *Pointer) AddButtonListener(channel chan interface{}) {
	p.addListener(3, channel)
}

func pointer_button(p *Pointer, msg []byte) {
	printEvent("button", msg)
	var data PointerButton
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

	button,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.button = button

	state,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.state = state

	for _,channel := range p.listeners[3] {
		channel <- data
	}
}

type PointerAxis struct {
	time uint32
	axis uint32
	value int32
}

func (p *Pointer) AddAxisListener(channel chan interface{}) {
	p.addListener(4, channel)
}

func pointer_axis(p *Pointer, msg []byte) {
	printEvent("axis", msg)
	var data PointerAxis
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.time = time

	axis,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.axis = axis

	value,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.value = value

	for _,channel := range p.listeners[4] {
		channel <- data
	}
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