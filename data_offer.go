package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataOffer struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *DataOffer, msg message)
}

//// Requests
func (d *DataOffer) Accept (serial uint32, typ string) {
	msg := newMessage(d, 0)
	writeInteger(msg,serial)
	writeString(msg,[]byte(typ))

	sendmsg(msg)
	printRequest("data_offer", d, "accept", serial, typ)
}

func (d *DataOffer) Receive (mime_type string, fd uintptr) {
	msg := newMessage(d, 1)
	writeString(msg,[]byte(mime_type))
	writeFd(msg,fd)

	sendmsg(msg)
	printRequest("data_offer", d, "receive", mime_type, fd)
}

func (d *DataOffer) Destroy () {
	msg := newMessage(d, 2)

	sendmsg(msg)
	printRequest("data_offer", d, "destroy")
}

//// Events
func (d *DataOffer) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

type DataOfferOffer struct {
	Typ string
}

func (d *DataOffer) AddOfferListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
	addListener(channel)
}

func data_offer_offer(d *DataOffer, msg message) {
	var data DataOfferOffer

	_,typ,err := readString(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Typ = typ

	for _,channel := range d.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("data_offer", d, "offer", typ)
}

func NewDataOffer() (d *DataOffer) {
	d = new(DataOffer)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, data_offer_offer)
	return
}

func (d *DataOffer) SetId(id int32) {
	d.id = id
}

func (d *DataOffer) Id() int32 {
	return d.id
}