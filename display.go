package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Display struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Display, msg message)
}

//// Requests
func (d *Display) Bind (name uint32, iface string, version uint32, id Object) {
	msg := newMessage(d, 0)
	writeInteger(msg,name)
	writeString(msg,[]byte(iface))
	writeInteger(msg,version)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("display", d, "bind", name, iface, version, "new id", id.Id())
}

func (d *Display) Sync (callback *Callback) {
	msg := newMessage(d, 1)
	appendObject(callback)
	writeInteger(msg,callback.Id())

	sendmsg(msg)
	printRequest("display", d, "sync", "new id", callback.Id())
}

//// Events
func (d *Display) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

type DisplayError struct {
	ObjectId Object
	Code uint32
	Message string
}

func (d *Display) AddErrorListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
	addListener(channel)
}

func display_error(d *Display, msg message) {
	var data DisplayError

	object_idid, err := readInt32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	var object_id Object
	object_idobj := getObject(object_idid)
	if object_idobj == nil {
		return
	}
	object_id = object_idobj.(Object)
	data.ObjectId = object_id

	code,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Code = code

	_,message,err := readString(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Message = message

	for _,channel := range d.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", d, "error", object_id.Id(), code, message)
}

type DisplayGlobal struct {
	Name uint32
	Iface string
	Version uint32
}

func (d *Display) AddGlobalListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
	addListener(channel)
}

func display_global(d *Display, msg message) {
	var data DisplayGlobal

	name,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Name = name

	_,iface,err := readString(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Iface = iface

	version,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Version = version

	for _,channel := range d.listeners[1] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", d, "global", name, iface, version)
}

type DisplayGlobalRemove struct {
	Name uint32
}

func (d *Display) AddGlobalRemoveListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
	addListener(channel)
}

func display_global_remove(d *Display, msg message) {
	var data DisplayGlobalRemove

	name,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Name = name

	for _,channel := range d.listeners[2] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", d, "global_remove", name)
}

type DisplayDeleteId struct {
	Id uint32
}

func (d *Display) AddDeleteIdListener(channel chan interface{}) {
	d.listeners[3] = append(d.listeners[3], channel)
	addListener(channel)
}

func display_delete_id(d *Display, msg message) {
	var data DisplayDeleteId

	id,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Id = id

	for _,channel := range d.listeners[3] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", d, "delete_id", id)
}

func NewDisplay() (d *Display) {
	d = new(Display)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, display_error)
	d.events = append(d.events, display_global)
	d.events = append(d.events, display_global_remove)
	d.events = append(d.events, display_delete_id)
	return
}

func (d *Display) SetId(id int32) {
	d.id = id
}

func (d *Display) Id() int32 {
	return d.id
}