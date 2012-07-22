package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Display struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Display, msg []byte)
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
	printRequest("display", "bind", name, iface, version, id)
}

func (d *Display) Sync (callback *Callback) {
	msg := newMessage(d, 1)
	appendObject(callback)
	writeInteger(msg,callback.Id())

	sendmsg(msg)
	printRequest("display", "sync", callback)
}

//// Events
func (d *Display) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type DisplayError struct {
	ObjectId Object
	Code uint32
	Message string
}

func (d *Display) AddErrorListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
}

func display_error(d *Display, msg []byte) {
	var data DisplayError
	buf := bytes.NewBuffer(msg)

	object_idid, err := readInt32(buf)
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

	code,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Code = code

	_,message,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Message = message

	for _,channel := range d.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", "error", object_id, code, message)
}

type DisplayGlobal struct {
	Name uint32
	Iface string
	Version uint32
}

func (d *Display) AddGlobalListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
}

func display_global(d *Display, msg []byte) {
	var data DisplayGlobal
	buf := bytes.NewBuffer(msg)

	name,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Name = name

	_,iface,err := readString(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Iface = iface

	version,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Version = version

	for _,channel := range d.listeners[1] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", "global", name, iface, version)
}

type DisplayGlobalRemove struct {
	Name uint32
}

func (d *Display) AddGlobalRemoveListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
}

func display_global_remove(d *Display, msg []byte) {
	var data DisplayGlobalRemove
	buf := bytes.NewBuffer(msg)

	name,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Name = name

	for _,channel := range d.listeners[2] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", "global_remove", name)
}

type DisplayDeleteId struct {
	Id uint32
}

func (d *Display) AddDeleteIdListener(channel chan interface{}) {
	d.listeners[3] = append(d.listeners[3], channel)
}

func display_delete_id(d *Display, msg []byte) {
	var data DisplayDeleteId
	buf := bytes.NewBuffer(msg)

	id,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Id = id

	for _,channel := range d.listeners[3] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("display", "delete_id", id)
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