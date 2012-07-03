
package gowl

import (
	"bytes"
)

type Display struct {
	*WlObject
	events []func (d *Display, msg []byte)
}

//// Requests
func (d *Display) Bind (name uint32, iface string, version uint32, id Object ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, name)
	writeString(buf, []byte(iface))
	writeInteger(buf, version)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(d, 0, buf.Bytes())
}

func (d *Display) Sync (callback *Callback ) {
	buf := new(bytes.Buffer)
	appendObject(callback)
	writeInteger(buf, callback.Id())

	sendmsg(d, 1, buf.Bytes())
}

//// Events
func (d *Display) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type DisplayError struct {
	object_id Object
	code uint32
	message string
}

func (d *Display) AddErrorListener(channel chan interface{}) {
	d.addListener(0, channel)
}

func display_error(d *Display, msg []byte) {
	var data DisplayError
	buf := bytes.NewBuffer(msg)

	object_idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	object_id := new(Object)
	object_id = getObject(object_idid).(Object)
	data.object_id = object_id

	code,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.code = code

	_,message,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.message = message

	for _,channel := range d.listeners[0] {
		channel <- data
	}
}

type DisplayGlobal struct {
	name uint32
	iface string
	version uint32
}

func (d *Display) AddGlobalListener(channel chan interface{}) {
	d.addListener(1, channel)
}

func display_global(d *Display, msg []byte) {
	var data DisplayGlobal
	buf := bytes.NewBuffer(msg)

	name,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.name = name

	_,iface,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.iface = iface

	version,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.version = version

	for _,channel := range d.listeners[1] {
		channel <- data
	}
}

type DisplayGlobal_remove struct {
	name uint32
}

func (d *Display) AddGlobal_removeListener(channel chan interface{}) {
	d.addListener(2, channel)
}

func display_global_remove(d *Display, msg []byte) {
	var data DisplayGlobal_remove
	buf := bytes.NewBuffer(msg)

	name,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.name = name

	for _,channel := range d.listeners[2] {
		channel <- data
	}
}

type DisplayDelete_id struct {
	id uint32
}

func (d *Display) AddDelete_idListener(channel chan interface{}) {
	d.addListener(3, channel)
}

func display_delete_id(d *Display, msg []byte) {
	var data DisplayDelete_id
	buf := bytes.NewBuffer(msg)

	id,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.id = id

	for _,channel := range d.listeners[3] {
		channel <- data
	}
}

func NewDisplay() (d *Display) {
	d = new(Display)

	d.events = append(d.events, display_error)
	d.events = append(d.events, display_global)
	d.events = append(d.events, display_global_remove)
	d.events = append(d.events, display_delete_id)
	return
}