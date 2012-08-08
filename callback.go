package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Callback struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (c *Callback, msg message)
}

//// Requests
//// Events
func (c *Callback) HandleEvent(msg message) {
	if c.events[msg.opcode] != nil {
		c.events[msg.opcode](c, msg)
	}
}

type CallbackDone struct {
	Serial uint32
}

func (c *Callback) AddDoneListener(channel chan interface{}) {
	c.listeners[0] = append(c.listeners[0], channel)
	addListener(channel)
}

func callback_done(c *Callback, msg message) {
	var data CallbackDone

	serial,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	for _,channel := range c.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("callback", c, "done", serial)
}

func NewCallback() (c *Callback) {
	c = new(Callback)
	c.listeners = make(map[int16][]chan interface{}, 0)

	c.events = append(c.events, callback_done)
	return
}

func (c *Callback) SetId(id int32) {
	c.id = id
}

func (c *Callback) Id() int32 {
	return c.id
}