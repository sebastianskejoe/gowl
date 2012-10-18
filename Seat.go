package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Seat struct {
	id int32
    capabilitiesListeners []chan SeatCapabilities
	events []func(s *Seat, msg message) error
    name string
}

func NewSeat() (s *Seat) {
	s = new(Seat)
    s.name = "Seat"
    s.capabilitiesListeners = make([]chan SeatCapabilities, 0)

    s.events = append(s.events, seatCapabilities)
	return
}

func (s *Seat) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

func (s *Seat) SetId(id int32) {
	s.id = id
}

func (s *Seat) Id() int32 {
	return s.id
}

func (s *Seat) Name() string {
    return s.name
}

////
//// REQUESTS
////

func (s *Seat) GetPointer(id *Pointer) {
    sendrequest(s, "wl_seat_get_pointer", id)
}

func (s *Seat) GetKeyboard(id *Keyboard) {
    sendrequest(s, "wl_seat_get_keyboard", id)
}

func (s *Seat) GetTouch(id *Touch) {
    sendrequest(s, "wl_seat_get_touch", id)
}

////
//// EVENTS
////

type SeatCapabilities struct {
    Capabilities uint32
}

func (s *Seat) AddCapabilitiesListener(channel chan SeatCapabilities) {
    s.capabilitiesListeners = append(s.capabilitiesListeners, channel)
}

func seatCapabilities(s *Seat, msg message) (err error) {
    var data SeatCapabilities

    // Read capabilities
    capabilities,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Capabilities = capabilities

    // Dispatch events
    for _,channel := range s.capabilitiesListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
