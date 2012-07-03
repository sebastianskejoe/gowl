
package gowl

import (
	"bytes"
)

type Output struct {
	*WlObject
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
	x int32
	y int32
	physical_width int32
	physical_height int32
	subpixel int32
	make string
	model string
}

func (o *Output) AddGeometryListener(channel chan interface{}) {
	o.addListener(0, channel)
}

func output_geometry(o *Output, msg []byte) {
	var data OutputGeometry
	buf := bytes.NewBuffer(msg)

	x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.x = x

	y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.y = y

	physical_width,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.physical_width = physical_width

	physical_height,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.physical_height = physical_height

	subpixel,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.subpixel = subpixel

	_,make,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.make = make

	_,model,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.model = model

	for _,channel := range o.listeners[0] {
		channel <- data
	}
}

type OutputMode struct {
	flags uint32
	width int32
	height int32
	refresh int32
}

func (o *Output) AddModeListener(channel chan interface{}) {
	o.addListener(1, channel)
}

func output_mode(o *Output, msg []byte) {
	var data OutputMode
	buf := bytes.NewBuffer(msg)

	flags,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.flags = flags

	width,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.width = width

	height,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.height = height

	refresh,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.refresh = refresh

	for _,channel := range o.listeners[1] {
		channel <- data
	}
}

func NewOutput() (o *Output) {
	o = new(Output)

	o.events = append(o.events, output_geometry)
	o.events = append(o.events, output_mode)
	return
}