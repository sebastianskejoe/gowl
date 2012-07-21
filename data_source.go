package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataSource struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *DataSource, msg []byte)
}

//// Requests
func (d *DataSource) Offer (typ string) {
	msg := newMessage(d, 0)
	writeString(msg,[]byte(typ))

	sendmsg(msg)
	printRequest("data_source", "offer", typ)
}

func (d *DataSource) Destroy () {
	msg := newMessage(d, 1)

	sendmsg(msg)
	printRequest("data_source", "destroy", )
}

//// Events
func (d *DataSource) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type DataSourceTarget struct {
	MimeType string
}

func (d *DataSource) AddTargetListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
}

func data_source_target(d *DataSource, msg []byte) {
	var data DataSourceTarget
	buf := bytes.NewBuffer(msg)

	_,mime_type,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.MimeType = mime_type

	for _,channel := range d.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("data_source", "target", mime_type)
}

type DataSourceSend struct {
	MimeType string
	Fd uintptr
}

func (d *DataSource) AddSendListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
}

func data_source_send(d *DataSource, msg []byte) {
	var data DataSourceSend
	buf := bytes.NewBuffer(msg)

	_,mime_type,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.MimeType = mime_type

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

type DataSourceCancelled struct {
}

func (d *DataSource) AddCancelledListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
}

func data_source_cancelled(d *DataSource, msg []byte) {
	var data DataSourceCancelled

	for _,channel := range d.listeners[2] {
		go func () { channel <- data }()
	}
	printEvent("data_source", "cancelled", )
}

func NewDataSource() (d *DataSource) {
	d = new(DataSource)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, data_source_target)
	d.events = append(d.events, data_source_send)
	d.events = append(d.events, data_source_cancelled)
	return
}

func (d *DataSource) SetId(id int32) {
	d.id = id
}

func (d *DataSource) Id() int32 {
	return d.id
}