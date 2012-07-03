
package gowl

import (
	"bytes"
)

type Shell_surface struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shell_surface, msg []byte)
}

//// Requests
func (s *Shell_surface) Pong (serial uint32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, serial)

	sendmsg(s, 0, buf.Bytes())
}

func (s *Shell_surface) Move (seat *Seat, serial uint32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, seat.Id())
	writeInteger(buf, serial)

	sendmsg(s, 1, buf.Bytes())
}

func (s *Shell_surface) Resize (seat *Seat, serial uint32, edges uint32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, seat.Id())
	writeInteger(buf, serial)
	writeInteger(buf, edges)

	sendmsg(s, 2, buf.Bytes())
}

func (s *Shell_surface) Set_toplevel ( ) {
	buf := new(bytes.Buffer)

	sendmsg(s, 3, buf.Bytes())
}

func (s *Shell_surface) Set_transient (parent *Surface, x int32, y int32, flags uint32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, parent.Id())
	writeInteger(buf, x)
	writeInteger(buf, y)
	writeInteger(buf, flags)

	sendmsg(s, 4, buf.Bytes())
}

func (s *Shell_surface) Set_fullscreen (method uint32, framerate uint32, output *Output ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, method)
	writeInteger(buf, framerate)
	writeInteger(buf, output.Id())

	sendmsg(s, 5, buf.Bytes())
}

func (s *Shell_surface) Set_popup (seat *Seat, serial uint32, parent *Surface, x int32, y int32, flags uint32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, seat.Id())
	writeInteger(buf, serial)
	writeInteger(buf, parent.Id())
	writeInteger(buf, x)
	writeInteger(buf, y)
	writeInteger(buf, flags)

	sendmsg(s, 6, buf.Bytes())
}

func (s *Shell_surface) Set_maximized (output *Output ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, output.Id())

	sendmsg(s, 7, buf.Bytes())
}

func (s *Shell_surface) Set_title (title string ) {
	buf := new(bytes.Buffer)
	writeString(buf, []byte(title))

	sendmsg(s, 8, buf.Bytes())
}

func (s *Shell_surface) Set_class (class_ string ) {
	buf := new(bytes.Buffer)
	writeString(buf, []byte(class_))

	sendmsg(s, 9, buf.Bytes())
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
	printEvent("ping", msg)
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
	printEvent("configure", msg)
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
}

type Shell_surfacePopup_done struct {
}

func (s *Shell_surface) AddPopup_doneListener(channel chan interface{}) {
	s.listeners[2] = append(s.listeners[2], channel)
}

func shell_surface_popup_done(s *Shell_surface, msg []byte) {
	printEvent("popup_done", msg)
	var data Shell_surfacePopup_done

	for _,channel := range s.listeners[2] {
		go func () { channel <- data }()
	}
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