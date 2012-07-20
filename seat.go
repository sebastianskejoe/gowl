package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Seat struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Seat, msg []byte)
}

//// Requests
func (s *Seat) GetPointer (id *Pointer) {
	msg := newMessage(s, 0)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("seat", "get_pointer", id)
}

func (s *Seat) GetKeyboard (id *Keyboard) {
	msg := newMessage(s, 1)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("seat", "get_keyboard", id)
}

func (s *Seat) GetTouch (id *Touch) {
	msg := newMessage(s, 2)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("seat", "get_touch", id)
}

//// Events
func (s *Seat) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

type SeatCapabilities struct {
	Capabilities uint32
}

func (s *Seat) AddCapabilitiesListener(channel chan interface{}) {
	s.listeners[0] = append(s.listeners[0], channel)
}

func seat_capabilities(s *Seat, msg []byte) {
	var data SeatCapabilities
	buf := bytes.NewBuffer(msg)

	capabilities,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Capabilities = capabilities

	for _,channel := range s.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("seat", "capabilities", capabilities)
}

func NewSeat() (s *Seat) {
	s = new(Seat)
	s.listeners = make(map[int16][]chan interface{}, 0)

	s.events = append(s.events, seat_capabilities)
	return
}

func (s *Seat) SetId(id int32) {
	s.id = id
}

func (s *Seat) Id() int32 {
	return s.id
}