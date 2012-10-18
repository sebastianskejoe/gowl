package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Display struct {
	id int32
    errorListeners []chan DisplayError
	deleteIdListeners []chan DisplayDeleteId
	events []func(d *Display, msg message) error
    name string
}

func NewDisplay() (d *Display) {
	d = new(Display)
    d.name = "Display"
    d.errorListeners = make([]chan DisplayError, 0)
	d.deleteIdListeners = make([]chan DisplayDeleteId, 0)

    d.events = append(d.events, displayError)
	d.events = append(d.events, displayDeleteId)
	return
}

func (d *Display) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

func (d *Display) SetId(id int32) {
	d.id = id
}

func (d *Display) Id() int32 {
	return d.id
}

func (d *Display) Name() string {
    return d.name
}

////
//// REQUESTS
////

func (d *Display) Sync(callback *Callback) {
    sendrequest(d, "wl_display_sync", callback)
}

func (d *Display) GetRegistry(callback *Registry) {
    sendrequest(d, "wl_display_get_registry", callback)
}

////
//// EVENTS
////

type DisplayError struct {
    ObjectId Object
	Code uint32
	Message string
}

func (d *Display) AddErrorListener(channel chan DisplayError) {
    d.errorListeners = append(d.errorListeners, channel)
}

func displayError(d *Display, msg message) (err error) {
    var data DisplayError

    // Read object_id
    object_id,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    object_idObj := getObject(object_id)
    data.ObjectId = object_idObj.(Object)

    // Read code
    code,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Code = code

    // Read message
    message,err := readString(msg.buf)
    if err != nil {
        return
    }
    data.Message = message

    // Dispatch events
    for _,channel := range d.errorListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DisplayDeleteId struct {
    Id uint32
}

func (d *Display) AddDeleteIdListener(channel chan DisplayDeleteId) {
    d.deleteIdListeners = append(d.deleteIdListeners, channel)
}

func displayDeleteId(d *Display, msg message) (err error) {
    var data DisplayDeleteId

    // Read id
    id,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Id = id

    // Dispatch events
    for _,channel := range d.deleteIdListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
