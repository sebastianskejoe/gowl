
package gowl

import (
	"bytes"
)

type Output struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (o *Output, msg []byte)
}

//// Requests
//// Events
func (o *Output) HandleEvent(opcode int16, msg []byte) {
	if o.events[opcode] != nil {
		o.events[opcode](o, msg)
	}
}

type OutputGeometry struct {
	X int32
	Y int32
	Physical_width int32
	Physical_height int32
	Subpixel int32
	Make string
	Model string
}

func (o *Output) AddGeometryListener(channel chan interface{}) {
	o.listeners[0] = append(o.listeners[0], channel)
}

func output_geometry(o *Output, msg []byte) {
	printEvent("geometry", msg)
	var data OutputGeometry
	buf := bytes.NewBuffer(msg)

	x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.X = x

	y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Y = y

	physical_width,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Physical_width = physical_width

	physical_height,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Physical_height = physical_height

	subpixel,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Subpixel = subpixel

	_,make,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Make = make

	_,model,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Model = model

	for _,channel := range o.listeners[0] {
		go func () { channel <- data }()
	}
}

type OutputMode struct {
	Flags uint32
	Width int32
	Height int32
	Refresh int32
}

func (o *Output) AddModeListener(channel chan interface{}) {
	o.listeners[1] = append(o.listeners[1], channel)
}

func output_mode(o *Output, msg []byte) {
	printEvent("mode", msg)
	var data OutputMode
	buf := bytes.NewBuffer(msg)

	flags,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Flags = flags

	width,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Width = width

	height,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Height = height

	refresh,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Refresh = refresh

	for _,channel := range o.listeners[1] {
		go func () { channel <- data }()
	}
}

func NewOutput() (o *Output) {
	o = new(Output)
	o.listeners = make(map[int16][]chan interface{}, 0)

	o.events = append(o.events, output_geometry)
	o.events = append(o.events, output_mode)
	return
}

func (o *Output) SetId(id int32) {
	o.id = id
}

func (o *Output) Id() int32 {
	return o.id
}