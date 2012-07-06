package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Data_offer struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Data_offer, msg []byte)
}

//// Requests
func (d *Data_offer) Accept (serial uint32, typ string) {
	msg := newMessage(d, 0)
	writeInteger(msg,serial)
	writeString(msg,[]byte(typ))

	sendmsg(msg)
	printRequest("data_offer", "accept", serial, typ)
}

func (d *Data_offer) Receive (mime_type string, fd uintptr) {
	msg := newMessage(d, 1)
	writeString(msg,[]byte(mime_type))
	writeFd(msg,fd)

	sendmsg(msg)
	printRequest("data_offer", "receive", mime_type, fd)
}

func (d *Data_offer) Destroy () {
	msg := newMessage(d, 2)

	sendmsg(msg)
	printRequest("data_offer", "destroy", )
}

//// Events
func (d *Data_offer) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type Data_offerOffer struct {
	Typ string
}

func (d *Data_offer) AddOfferListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
}

func data_offer_offer(d *Data_offer, msg []byte) {
	var data Data_offerOffer
	buf := bytes.NewBuffer(msg)

	_,typ,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Typ = typ

	for _,channel := range d.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("data_offer", "offer", typ)
}

func NewData_offer() (d *Data_offer) {
	d = new(Data_offer)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, data_offer_offer)
	return
}

func (d *Data_offer) SetId(id int32) {
	d.id = id
}

func (d *Data_offer) Id() int32 {
	return d.id
}