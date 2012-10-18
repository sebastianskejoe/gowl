package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Callback struct {
	id int32
    doneListeners []chan CallbackDone
	events []func(c *Callback, msg message) error
    name string
}

func NewCallback() (c *Callback) {
	c = new(Callback)
    c.name = "Callback"
    c.doneListeners = make([]chan CallbackDone, 0)

    c.events = append(c.events, callbackDone)
	return
}

func (c *Callback) HandleEvent(msg message) {
	if c.events[msg.opcode] != nil {
		c.events[msg.opcode](c, msg)
	}
}

func (c *Callback) SetId(id int32) {
	c.id = id
}

func (c *Callback) Id() int32 {
	return c.id
}

func (c *Callback) Name() string {
    return c.name
}

////
//// REQUESTS
////

////
//// EVENTS
////

type CallbackDone struct {
    Serial uint32
}

func (c *Callback) AddDoneListener(channel chan CallbackDone) {
    c.doneListeners = append(c.doneListeners, channel)
}

func callbackDone(c *Callback, msg message) (err error) {
    var data CallbackDone

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Dispatch events
    for _,channel := range c.doneListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
