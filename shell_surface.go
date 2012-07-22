package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type ShellSurface struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *ShellSurface, msg []byte)
}

//// Requests
func (s *ShellSurface) Pong (serial uint32) {
	msg := newMessage(s, 0)
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("shell_surface", s, "pong", serial)
}

func (s *ShellSurface) Move (seat *Seat, serial uint32) {
	msg := newMessage(s, 1)
	writeInteger(msg,seat.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("shell_surface", s, "move", seat.Id(), serial)
}

func (s *ShellSurface) Resize (seat *Seat, serial uint32, edges uint32) {
	msg := newMessage(s, 2)
	writeInteger(msg,seat.Id())
	writeInteger(msg,serial)
	writeInteger(msg,edges)

	sendmsg(msg)
	printRequest("shell_surface", s, "resize", seat.Id(), serial, edges)
}

func (s *ShellSurface) SetToplevel () {
	msg := newMessage(s, 3)

	sendmsg(msg)
	printRequest("shell_surface", s, "set_toplevel")
}

func (s *ShellSurface) SetTransient (parent *Surface, x int32, y int32, flags uint32) {
	msg := newMessage(s, 4)
	writeInteger(msg,parent.Id())
	writeInteger(msg,x)
	writeInteger(msg,y)
	writeInteger(msg,flags)

	sendmsg(msg)
	printRequest("shell_surface", s, "set_transient", parent.Id(), x, y, flags)
}

func (s *ShellSurface) SetFullscreen (method uint32, framerate uint32, output *Output) {
	msg := newMessage(s, 5)
	writeInteger(msg,method)
	writeInteger(msg,framerate)
	writeInteger(msg,output.Id())

	sendmsg(msg)
	printRequest("shell_surface", s, "set_fullscreen", method, framerate, output.Id())
}

func (s *ShellSurface) SetPopup (seat *Seat, serial uint32, parent *Surface, x int32, y int32, flags uint32) {
	msg := newMessage(s, 6)
	writeInteger(msg,seat.Id())
	writeInteger(msg,serial)
	writeInteger(msg,parent.Id())
	writeInteger(msg,x)
	writeInteger(msg,y)
	writeInteger(msg,flags)

	sendmsg(msg)
	printRequest("shell_surface", s, "set_popup", seat.Id(), serial, parent.Id(), x, y, flags)
}

func (s *ShellSurface) SetMaximized (output *Output) {
	msg := newMessage(s, 7)
	writeInteger(msg,output.Id())

	sendmsg(msg)
	printRequest("shell_surface", s, "set_maximized", output.Id())
}

func (s *ShellSurface) SetTitle (title string) {
	msg := newMessage(s, 8)
	writeString(msg,[]byte(title))

	sendmsg(msg)
	printRequest("shell_surface", s, "set_title", title)
}

func (s *ShellSurface) SetClass (class_ string) {
	msg := newMessage(s, 9)
	writeString(msg,[]byte(class_))

	sendmsg(msg)
	printRequest("shell_surface", s, "set_class", class_)
}

//// Events
func (s *ShellSurface) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

type ShellSurfacePing struct {
	Serial uint32
}

func (s *ShellSurface) AddPingListener(channel chan interface{}) {
	s.listeners[0] = append(s.listeners[0], channel)
}

func shell_surface_ping(s *ShellSurface, msg []byte) {
	var data ShellSurfacePing
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	for _,channel := range s.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("shell_surface", s, "ping", serial)
}

type ShellSurfaceConfigure struct {
	Edges uint32
	Width int32
	Height int32
}

func (s *ShellSurface) AddConfigureListener(channel chan interface{}) {
	s.listeners[1] = append(s.listeners[1], channel)
}

func shell_surface_configure(s *ShellSurface, msg []byte) {
	var data ShellSurfaceConfigure
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
		go func() {
			channel <- data
		} ()
	}
	printEvent("shell_surface", s, "configure", edges, width, height)
}

type ShellSurfacePopupDone struct {
}

func (s *ShellSurface) AddPopupDoneListener(channel chan interface{}) {
	s.listeners[2] = append(s.listeners[2], channel)
}

func shell_surface_popup_done(s *ShellSurface, msg []byte) {
	var data ShellSurfacePopupDone

	for _,channel := range s.listeners[2] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("shell_surface", s, "popup_done")
}

func NewShellSurface() (s *ShellSurface) {
	s = new(ShellSurface)
	s.listeners = make(map[int16][]chan interface{}, 0)

	s.events = append(s.events, shell_surface_ping)
	s.events = append(s.events, shell_surface_configure)
	s.events = append(s.events, shell_surface_popup_done)
	return
}

func (s *ShellSurface) SetId(id int32) {
	s.id = id
}

func (s *ShellSurface) Id() int32 {
	return s.id
}