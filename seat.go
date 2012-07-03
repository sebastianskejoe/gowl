
package gowl

import (
	"bytes"
)

type Seat struct {
	*WlObject
	events []func (s *Seat, msg []byte)
}

//// Requests
func (s *Seat) Get_pointer (id *Pointer ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(s, 0, buf.Bytes())
}

func (s *Seat) Get_keyboard (id *Keyboard ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(s, 1, buf.Bytes())
}

func (s *Seat) Get_touch (id *Touch ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(s, 2, buf.Bytes())
}

//// Events
func (s *Seat) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

type SeatCapabilities struct {
	capabilities uint32
}

func (s *Seat) AddCapabilitiesListener(channel chan interface{}) {
	s.addListener(0, channel)
}

func seat_capabilities(s *Seat, msg []byte) {
	printEvent("capabilities", msg)
	var data SeatCapabilities
	buf := bytes.NewBuffer(msg)

	capabilities,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.capabilities = capabilities

	for _,channel := range s.listeners[0] {
		channel <- data
	}
}

func NewSeat() (s *Seat) {
	s = new(Seat)
	s.listeners = make(map[int16]chan interface{}, 0)

	s.events = append(s.events, seat_capabilities)
	return
}