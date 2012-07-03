
package gowl

import (
	"bytes"
)

type Data_offer struct {
	*WlObject
	events []func (d *Data_offer, msg []byte)
}

//// Requests
func (d *Data_offer) Accept (serial uint32, typ string ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, serial)
	writeString(buf, []byte(typ))

	sendmsg(d, 0, buf.Bytes())
}

func (d *Data_offer) Receive (mime_type string, fd uintptr ) {
	buf := new(bytes.Buffer)
	writeString(buf, []byte(mime_type))
	writeInteger(buf, fd)

	sendmsg(d, 1, buf.Bytes())
}

func (d *Data_offer) Destroy ( ) {
	buf := new(bytes.Buffer)

	sendmsg(d, 2, buf.Bytes())
}

//// Events
func (d *Data_offer) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type Data_offerOffer struct {
	typ string
}

func (d *Data_offer) AddOfferListener(channel chan interface{}) {
	d.addListener(0, channel)
}

func data_offer_offer(d *Data_offer, msg []byte) {
	printEvent("offer", msg)
	var data Data_offerOffer
	buf := bytes.NewBuffer(msg)

	_,typ,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.typ = typ

	for _,channel := range d.listeners[0] {
		channel <- data
	}
}

func NewData_offer() (d *Data_offer) {
	d = new(Data_offer)
	d.listeners = make(map[int16]chan interface{}, 0)

	d.events = append(d.events, data_offer_offer)
	return
}