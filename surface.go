package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Surface struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Surface, msg []byte)
}

//// Requests
func (s *Surface) Destroy () {
	msg := newMessage(s, 0)

	sendmsg(msg)
	printRequest("surface", "destroy", )
}

func (s *Surface) Attach (buffer *Buffer, x int32, y int32) {
	msg := newMessage(s, 1)
	writeInteger(msg,buffer.Id())
	writeInteger(msg,x)
	writeInteger(msg,y)

	sendmsg(msg)
	printRequest("surface", "attach", buffer, x, y)
}

func (s *Surface) Damage (x int32, y int32, width int32, height int32) {
	msg := newMessage(s, 2)
	writeInteger(msg,x)
	writeInteger(msg,y)
	writeInteger(msg,width)
	writeInteger(msg,height)

	sendmsg(msg)
	printRequest("surface", "damage", x, y, width, height)
}

func (s *Surface) Frame (callback *Callback) {
	msg := newMessage(s, 3)
	appendObject(callback)
	writeInteger(msg,callback.Id())

	sendmsg(msg)
	printRequest("surface", "frame", callback)
}

func (s *Surface) SetOpaqueRegion (region *Region) {
	msg := newMessage(s, 4)
	writeInteger(msg,region.Id())

	sendmsg(msg)
	printRequest("surface", "set_opaque_region", region)
}

func (s *Surface) SetInputRegion (region *Region) {
	msg := newMessage(s, 5)
	writeInteger(msg,region.Id())

	sendmsg(msg)
	printRequest("surface", "set_input_region", region)
}

//// Events
func (s *Surface) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

type SurfaceEnter struct {
	Output *Output
}

func (s *Surface) AddEnterListener(channel chan interface{}) {
	s.listeners[0] = append(s.listeners[0], channel)
}

func surface_enter(s *Surface, msg []byte) {
	var data SurfaceEnter
	buf := bytes.NewBuffer(msg)

	outputid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	output := new(Output)
	outputobj := getObject(outputid)
	if outputobj == nil {
		return
	}
	output = outputobj.(*Output)
	data.Output = output

	for _,channel := range s.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("surface", "enter", output)
}

type SurfaceLeave struct {
	Output *Output
}

func (s *Surface) AddLeaveListener(channel chan interface{}) {
	s.listeners[1] = append(s.listeners[1], channel)
}

func surface_leave(s *Surface, msg []byte) {
	var data SurfaceLeave
	buf := bytes.NewBuffer(msg)

	outputid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	output := new(Output)
	outputobj := getObject(outputid)
	if outputobj == nil {
		return
	}
	output = outputobj.(*Output)
	data.Output = output

	for _,channel := range s.listeners[1] {
		go func () { channel <- data }()
	}
	printEvent("surface", "leave", output)
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