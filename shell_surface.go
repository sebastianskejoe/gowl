package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Shell_surface struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shell_surface, msg []byte)
}

//// Requests
func (s *Shell_surface) Pong (serial uint32) {
	msg := newMessage(s, 0)
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("shell_surface", "pong", serial)
}

func (s *Shell_surface) Move (seat *Seat, serial uint32) {
	msg := newMessage(s, 1)
	writeInteger(msg,seat.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("shell_surface", "move", seat, serial)
}

func (s *Shell_surface) Resize (seat *Seat, serial uint32, edges uint32) {
	msg := newMessage(s, 2)
	writeInteger(msg,seat.Id())
	writeInteger(msg,serial)
	writeInteger(msg,edges)

	sendmsg(msg)
	printRequest("shell_surface", "resize", seat, serial, edges)
}

func (s *Shell_surface) Set_toplevel () {
	msg := newMessage(s, 3)

	sendmsg(msg)
	printRequest("shell_surface", "set_toplevel", )
}

func (s *Shell_surface) Set_transient (parent *Surface, x int32, y int32, flags uint32) {
	msg := newMessage(s, 4)
	writeInteger(msg,parent.Id())
	writeInteger(msg,x)
	writeInteger(msg,y)
	writeInteger(msg,flags)

	sendmsg(msg)
	printRequest("shell_surface", "set_transient", parent, x, y, flags)
}

func (s *Shell_surface) Set_fullscreen (method uint32, framerate uint32, output *Output) {
	msg := newMessage(s, 5)
	writeInteger(msg,method)
	writeInteger(msg,framerate)
	writeInteger(msg,output.Id())

	sendmsg(msg)
	printRequest("shell_surface", "set_fullscreen", method, framerate, output)
}

func (s *Shell_surface) Set_popup (seat *Seat, serial uint32, parent *Surface, x int32, y int32, flags uint32) {
	msg := newMessage(s, 6)
	writeInteger(msg,seat.Id())
	writeInteger(msg,serial)
	writeInteger(msg,parent.Id())
	writeInteger(msg,x)
	writeInteger(msg,y)
	writeInteger(msg,flags)

	sendmsg(msg)
	printRequest("shell_surface", "set_popup", seat, serial, parent, x, y, flags)
}

func (s *Shell_surface) Set_maximized (output *Output) {
	msg := newMessage(s, 7)
	writeInteger(msg,output.Id())

	sendmsg(msg)
	printRequest("shell_surface", "set_maximized", output)
}

func (s *Shell_surface) Set_title (title string) {
	msg := newMessage(s, 8)
	writeString(msg,[]byte(title))

	sendmsg(msg)
	printRequest("shell_surface", "set_title", title)
}

func (s *Shell_surface) Set_class (class_ string) {
	msg := newMessage(s, 9)
	writeString(msg,[]byte(class_))

	sendmsg(msg)
	printRequest("shell_surface", "set_class", class_)
}

//// Events
func (s *Shell_surface) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

type Shell_surfacePing struct {
	Serial uint32
}

func (s *Shell_surface) AddPingListener(channel chan interface{}) {
	s.listeners[0] = append(s.listeners[0], channel)
}

func shell_surface_ping(s *Shell_surface, msg []byte) {
	var data Shell_surfacePing
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	for _,channel := range s.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("shell_surface", "ping", serial)
}

type Shell_surfaceConfigure struct {
	Edges uint32
	Width int32
	Height int32
}

func (s *Shell_surface) AddConfigureListener(channel chan interface{}) {
	s.listeners[1] = append(s.listeners[1], channel)
}

func shell_surface_configure(s *Shell_surface, msg []byte) {
	var data Shell_surfaceConfigure
	buf := bytes.NewBuffer(msg)

	edges,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Edges = edges

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

	for _,channel := range s.listeners[1] {
		go func () { channel <- data }()
	}
	printEvent("shell_surface", "configure", edges, width, height)
}

type Shell_surfacePopup_done struct {
}

func (s *Shell_surface) AddPopup_doneListener(channel chan interface{}) {
	s.listeners[2] = append(s.listeners[2], channel)
}

func shell_surface_popup_done(s *Shell_surface, msg []byte) {
	var data Shell_surfacePopup_done

	for _,channel := range s.listeners[2] {
		go func () { channel <- data }()
	}
	printEvent("shell_surface", "popup_done", )
}

func NewShell_surface() (s *Shell_surface) {
	s = new(Shell_surface)
	s.listeners = make(map[int16][]chan interface{}, 0)

	s.events = append(s.events, shell_surface_ping)
	s.events = append(s.events, shell_surface_configure)
	s.events = append(s.events, shell_surface_popup_done)
	return
}

func (s *Shell_surface) SetId(id int32) {
	s.id = id
}

func (s *Shell_surface) Id() int32 {
	return s.id
}