
package gowl

import (
	"bytes"
)

type Callback struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (c *Callback, msg []byte)
}

//// Requests
//// Events
func (c *Callback) HandleEvent(opcode int16, msg []byte) {
	if c.events[opcode] != nil {
		c.events[opcode](c, msg)
	}
}

type CallbackDone struct {
	Serial uint32
}

func (c *Callback) AddDoneListener(channel chan interface{}) {
	c.listeners[0] = append(c.listeners[0], channel)
}

func callback_done(c *Callback, msg []byte) {
	printEvent("done", msg)
	var data CallbackDone
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	for _,channel := range c.listeners[0] {
		go func () { channel <- data }()
	}
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