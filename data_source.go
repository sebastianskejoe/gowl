package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Data_source struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Data_source, msg []byte)
}

//// Requests
func (d *Data_source) Offer (typ string) {
	msg := newMessage(d, 0)
	writeString(msg,[]byte(typ))

	sendmsg(msg)
	printRequest("data_source", "offer", typ)
}

func (d *Data_source) Destroy () {
	msg := newMessage(d, 1)

	sendmsg(msg)
	printRequest("data_source", "destroy", )
}

//// Events
func (d *Data_source) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type Data_sourceTarget struct {
	Mime_type string
}

func (d *Data_source) AddTargetListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
}

func data_source_target(d *Data_source, msg []byte) {
	var data Data_sourceTarget
	buf := bytes.NewBuffer(msg)

	_,mime_type,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Mime_type = mime_type

	for _,channel := range d.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("data_source", "target", mime_type)
}

type Data_sourceSend struct {
	Mime_type string
	Fd uintptr
}

func (d *Data_source) AddSendListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
}

func data_source_send(d *Data_source, msg []byte) {
	var data Data_sourceSend
	buf := bytes.NewBuffer(msg)

	_,mime_type,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Mime_type = mime_type

	fd,err := readUintptr(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Fd = fd

	for _,channel := range d.listeners[1] {
		go func () { channel <- data }()
	}
	printEvent("data_source", "send", mime_type, fd)
}

type Data_sourceCancelled struct {
}

func (d *Data_source) AddCancelledListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
}

func data_source_cancelled(d *Data_source, msg []byte) {
	var data Data_sourceCancelled

	for _,channel := range d.listeners[2] {
		go func () { channel <- data }()
	}
	printEvent("data_source", "cancelled", )
}

func NewData_source() (d *Data_source) {
	d = new(Data_source)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, data_source_target)
	d.events = append(d.events, data_source_send)
	d.events = append(d.events, data_source_cancelled)
	return
}

func (d *Data_source) SetId(id int32) {
	d.id = id
}

func (d *Data_source) Id() int32 {
	return d.id
}