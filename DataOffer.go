package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataOffer struct {
	id int32
    offerListeners []chan DataOfferOffer
	events []func(d *DataOffer, msg message) error
    name string
}

func NewDataOffer() (d *DataOffer) {
	d = new(DataOffer)
    d.name = "DataOffer"
    d.offerListeners = make([]chan DataOfferOffer, 0)

    d.events = append(d.events, dataOfferOffer)
	return
}

func (d *DataOffer) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

func (d *DataOffer) SetId(id int32) {
	d.id = id
}

func (d *DataOffer) Id() int32 {
	return d.id
}

func (d *DataOffer) Name() string {
    return d.name
}

////
//// REQUESTS
////

func (d *DataOffer) Accept(serial uint32, typ string) {
    sendrequest(d, "wl_data_offer_accept", serial, typ)
}

func (d *DataOffer) Receive(mime_type string, fd uintptr) {
    sendrequest(d, "wl_data_offer_receive", mime_type, fd)
}

func (d *DataOffer) Destroy() {
    sendrequest(d, "wl_data_offer_destroy", )
}

////
//// EVENTS
////

type DataOfferOffer struct {
    Typ string
}

func (d *DataOffer) AddOfferListener(channel chan DataOfferOffer) {
    d.offerListeners = append(d.offerListeners, channel)
}

func dataOfferOffer(d *DataOffer, msg message) (err error) {
    var data DataOfferOffer

    // Read typ
    typ,err := readString(msg.buf)
    if err != nil {
        return
    }
    data.Typ = typ

    // Dispatch events
    for _,channel := range d.offerListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
