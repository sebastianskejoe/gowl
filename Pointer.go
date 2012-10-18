package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Pointer struct {
	id int32
    enterListeners []chan PointerEnter
	leaveListeners []chan PointerLeave
	motionListeners []chan PointerMotion
	buttonListeners []chan PointerButton
	axisListeners []chan PointerAxis
	events []func(p *Pointer, msg message) error
    name string
}

func NewPointer() (p *Pointer) {
	p = new(Pointer)
    p.name = "Pointer"
    p.enterListeners = make([]chan PointerEnter, 0)
	p.leaveListeners = make([]chan PointerLeave, 0)
	p.motionListeners = make([]chan PointerMotion, 0)
	p.buttonListeners = make([]chan PointerButton, 0)
	p.axisListeners = make([]chan PointerAxis, 0)

    p.events = append(p.events, pointerEnter)
	p.events = append(p.events, pointerLeave)
	p.events = append(p.events, pointerMotion)
	p.events = append(p.events, pointerButton)
	p.events = append(p.events, pointerAxis)
	return
}

func (p *Pointer) HandleEvent(msg message) {
	if p.events[msg.opcode] != nil {
		p.events[msg.opcode](p, msg)
	}
}

func (p *Pointer) SetId(id int32) {
	p.id = id
}

func (p *Pointer) Id() int32 {
	return p.id
}

func (p *Pointer) Name() string {
    return p.name
}

////
//// REQUESTS
////

func (p *Pointer) SetCursor(serial uint32, surface *Surface, hotspot_x int32, hotspot_y int32) {
    sendrequest(p, "wl_pointer_set_cursor", serial, surface, hotspot_x, hotspot_y)
}

////
//// EVENTS
////

type PointerEnter struct {
    Serial uint32
	Surface *Surface
	SurfaceX int32
	SurfaceY int32
}

func (p *Pointer) AddEnterListener(channel chan PointerEnter) {
    p.enterListeners = append(p.enterListeners, channel)
}

func pointerEnter(p *Pointer, msg message) (err error) {
    var data PointerEnter

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

    // Read surface_x
    surface_x,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.SurfaceX = surface_x

    // Read surface_y
    surface_y,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.SurfaceY = surface_y

    // Dispatch events
    for _,channel := range p.enterListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type PointerLeave struct {
    Serial uint32
	Surface *Surface
}

func (p *Pointer) AddLeaveListener(channel chan PointerLeave) {
    p.leaveListeners = append(p.leaveListeners, channel)
}

func pointerLeave(p *Pointer, msg message) (err error) {
    var data PointerLeave

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
    for _,channel := range p.leaveListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type PointerMotion struct {
    Time uint32
	SurfaceX int32
	SurfaceY int32
}

func (p *Pointer) AddMotionListener(channel chan PointerMotion) {
    p.motionListeners = append(p.motionListeners, channel)
}

func pointerMotion(p *Pointer, msg message) (err error) {
    var data PointerMotion

    // Read time
    time,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Time = time

    // Read surface_x
    surface_x,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.SurfaceX = surface_x

    // Read surface_y
    surface_y,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.SurfaceY = surface_y

    // Dispatch events
    for _,channel := range p.motionListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type PointerButton struct {
    Serial uint32
	Time uint32
	Button uint32
	State uint32
}

func (p *Pointer) AddButtonListener(channel chan PointerButton) {
    p.buttonListeners = append(p.buttonListeners, channel)
}

func pointerButton(p *Pointer, msg message) (err error) {
    var data PointerButton

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

    // Read button
    button,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Button = button

    // Read state
    state,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.State = state

    // Dispatch events
    for _,channel := range p.buttonListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type PointerAxis struct {
    Time uint32
	Axis uint32
	Value int32
}

func (p *Pointer) AddAxisListener(channel chan PointerAxis) {
    p.axisListeners = append(p.axisListeners, channel)
}

func pointerAxis(p *Pointer, msg message) (err error) {
    var data PointerAxis

    // Read time
    time,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Time = time

    // Read axis
    axis,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Axis = axis

    // Read value
    value,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.Value = value

    // Dispatch events
    for _,channel := range p.axisListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
