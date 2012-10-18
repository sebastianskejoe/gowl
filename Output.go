package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Output struct {
	id int32
    geometryListeners []chan OutputGeometry
	modeListeners []chan OutputMode
	events []func(o *Output, msg message) error
    name string
}

func NewOutput() (o *Output) {
	o = new(Output)
    o.name = "Output"
    o.geometryListeners = make([]chan OutputGeometry, 0)
	o.modeListeners = make([]chan OutputMode, 0)

    o.events = append(o.events, outputGeometry)
	o.events = append(o.events, outputMode)
	return
}

func (o *Output) HandleEvent(msg message) {
	if o.events[msg.opcode] != nil {
		o.events[msg.opcode](o, msg)
	}
}

func (o *Output) SetId(id int32) {
	o.id = id
}

func (o *Output) Id() int32 {
	return o.id
}

func (o *Output) Name() string {
    return o.name
}

////
//// REQUESTS
////

////
//// EVENTS
////

type OutputGeometry struct {
    X int32
	Y int32
	PhysicalWidth int32
	PhysicalHeight int32
	Subpixel int32
	Make string
	Model string
	Transform int32
}

func (o *Output) AddGeometryListener(channel chan OutputGeometry) {
    o.geometryListeners = append(o.geometryListeners, channel)
}

func outputGeometry(o *Output, msg message) (err error) {
    var data OutputGeometry

    // Read x
    x,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.X = x

    // Read y
    y,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Y = y

    // Read physical_width
    physical_width,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.PhysicalWidth = physical_width

    // Read physical_height
    physical_height,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.PhysicalHeight = physical_height

    // Read subpixel
    subpixel,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Subpixel = subpixel

    // Read make
    make,err := readString(msg.buf)
    if err != nil {
        return
    }
    data.Make = make

    // Read model
    model,err := readString(msg.buf)
    if err != nil {
        return
    }
    data.Model = model

    // Read transform
    transform,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Transform = transform

    // Dispatch events
    for _,channel := range o.geometryListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type OutputMode struct {
    Flags uint32
	Width int32
	Height int32
	Refresh int32
}

func (o *Output) AddModeListener(channel chan OutputMode) {
    o.modeListeners = append(o.modeListeners, channel)
}

func outputMode(o *Output, msg message) (err error) {
    var data OutputMode

    // Read flags
    flags,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Flags = flags

    // Read width
    width,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Width = width

    // Read height
    height,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Height = height

    // Read refresh
    refresh,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Refresh = refresh

    // Dispatch events
    for _,channel := range o.modeListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
