
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
	Serial uint32
	Surface *Surface
	Surface_x int32
	Surface_y int32
}

func (p *Pointer) AddEnterListener(channel chan interface{}) {
	p.listeners[0] = append(p.listeners[0], channel)
}

func pointer_enter(p *Pointer, msg []byte) {
	printEvent("enter", msg)
	var data PointerEnter
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

	surface_x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Surface_x = surface_x

	surface_y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Surface_y = surface_y

	for _,channel := range p.listeners[0] {
		go func () { channel <- data }()
	}
}

type PointerLeave struct {
	Serial uint32
	Surface *Surface
}

func (p *Pointer) AddLeaveListener(channel chan interface{}) {
	p.listeners[1] = append(p.listeners[1], channel)
}

func pointer_leave(p *Pointer, msg []byte) {
	printEvent("leave", msg)
	var data PointerLeave
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

	for _,channel := range p.listeners[1] {
		go func () { channel <- data }()
	}
}

type PointerMotion struct {
	Time uint32
	Surface_x int32
	Surface_y int32
}

func (p *Pointer) AddMotionListener(channel chan interface{}) {
	p.listeners[2] = append(p.listeners[2], channel)
}

func pointer_motion(p *Pointer, msg []byte) {
	printEvent("motion", msg)
	var data PointerMotion
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	surface_x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Surface_x = surface_x

	surface_y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Surface_y = surface_y

	for _,channel := range p.listeners[2] {
		go func () { channel <- data }()
	}
}

type PointerButton struct {
	Serial uint32
	Time uint32
	Button uint32
	State uint32
}

func (p *Pointer) AddButtonListener(channel chan interface{}) {
	p.listeners[3] = append(p.listeners[3], channel)
}

func pointer_button(p *Pointer, msg []byte) {
	printEvent("button", msg)
	var data PointerButton
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

	button,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Button = button

	state,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.State = state

	for _,channel := range p.listeners[3] {
		go func () { channel <- data }()
	}
}

type PointerAxis struct {
	Time uint32
	Axis uint32
	Value int32
}

func (p *Pointer) AddAxisListener(channel chan interface{}) {
	p.listeners[4] = append(p.listeners[4], channel)
}

func pointer_axis(p *Pointer, msg []byte) {
	printEvent("axis", msg)
	var data PointerAxis
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	axis,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Axis = axis

	value,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Value = value

	for _,channel := range p.listeners[4] {
		go func () { channel <- data }()
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

func (p *Pointer) SetId(id int32) {
	p.id = id
}

func (p *Pointer) Id() int32 {
	return p.id
}