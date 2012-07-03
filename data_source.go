
package gowl

import (
	"bytes"
)

type Data_source struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Data_source, msg []byte)
}

//// Requests
func (d *Data_source) Offer (typ string ) {
	buf := new(bytes.Buffer)
	writeString(buf, []byte(typ))

	sendmsg(d, 0, buf.Bytes())
}

func (d *Data_source) Destroy ( ) {
	buf := new(bytes.Buffer)

	sendmsg(d, 1, buf.Bytes())
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
	printEvent("target", msg)
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
}

type Data_sourceSend struct {
	Mime_type string
	Fd uintptr
}

func (d *Data_source) AddSendListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
}

func data_source_send(d *Data_source, msg []byte) {
	printEvent("send", msg)
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
}

type Data_sourceCancelled struct {
}

func (d *Data_source) AddCancelledListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
}

func data_source_cancelled(d *Data_source, msg []byte) {
	printEvent("cancelled", msg)
	var data Data_sourceCancelled

	for _,channel := range d.listeners[2] {
		go func () { channel <- data }()
	}
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