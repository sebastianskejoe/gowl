
package gowl

import (
	"bytes"
)

type Surface struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Surface, msg []byte)
}

//// Requests
func (s *Surface) Destroy ( ) {
	buf := new(bytes.Buffer)

	sendmsg(s, 0, buf.Bytes())
}

func (s *Surface) Attach (buffer *Buffer, x int32, y int32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, buffer.Id())
	writeInteger(buf, x)
	writeInteger(buf, y)

	sendmsg(s, 1, buf.Bytes())
}

func (s *Surface) Damage (x int32, y int32, width int32, height int32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, x)
	writeInteger(buf, y)
	writeInteger(buf, width)
	writeInteger(buf, height)

	sendmsg(s, 2, buf.Bytes())
}

func (s *Surface) Frame (callback *Callback ) {
	buf := new(bytes.Buffer)
	appendObject(callback)
	writeInteger(buf, callback.Id())

	sendmsg(s, 3, buf.Bytes())
}

func (s *Surface) Set_opaque_region (region *Region ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, region.Id())

	sendmsg(s, 4, buf.Bytes())
}

func (s *Surface) Set_input_region (region *Region ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, region.Id())

	sendmsg(s, 5, buf.Bytes())
}

//// Events
func (s *Surface) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

type SurfaceEnter struct {
	output *Output
}

func (s *Surface) AddEnterListener(channel chan interface{}) {
	s.listeners[0] = append(s.listeners[0], channel)
}

func surface_enter(s *Surface, msg []byte) {
	printEvent("enter", msg)
	var data SurfaceEnter
	buf := bytes.NewBuffer(msg)

	outputid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	output := new(Output)
	output = getObject(outputid).(*Output)
	data.output = output

	for _,channel := range s.listeners[0] {
		channel <- data
	}
}

type SurfaceLeave struct {
	output *Output
}

func (s *Surface) AddLeaveListener(channel chan interface{}) {
	s.listeners[1] = append(s.listeners[1], channel)
}

func surface_leave(s *Surface, msg []byte) {
	printEvent("leave", msg)
	var data SurfaceLeave
	buf := bytes.NewBuffer(msg)

	outputid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	output := new(Output)
	output = getObject(outputid).(*Output)
	data.output = output

	for _,channel := range s.listeners[1] {
		channel <- data
	}
}

func NewSurface() (s *Surface) {
	s = new(Surface)
	s.listeners = make(map[int16][]chan interface{}, 0)

	s.events = append(s.events, surface_enter)
	s.events = append(s.events, surface_leave)
	return
}

func (s *Surface) SetId(id int32) {
	s.id = id
}

func (s *Surface) Id() int32 {
	return s.id
}