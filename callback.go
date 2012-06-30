package gowl

import (
	"fmt"
	"bytes"
)

type Callback struct {
	events map[int16]func()
	id int32
	done chan uint32
}

func NewCallback(ch chan uint32) *Callback {
	cb := new(Callback)
	cb.events = make(map[int16]func())
	cb.done = ch
	return cb
}

func (c *Callback) HandleEvent(opcode int16, msg []byte) {
	if opcode != 0 {
		fmt.Println("Unknown opcode in callback", opcode)
		return
	}
	callback_done(c, msg)
	return
}

func (c *Callback) SetID(id int32) {
	c.id = id
}

func (c *Callback) ID() int32 {
	return c.id
}

func callback_done(cb *Callback, msg []byte) {
	serial,err := readUint32(bytes.NewBuffer(msg))
	if err != nil {
		printError("callback_done", err)
		return
	}
	go func () {cb.done <- serial}()
	fmt.Println("CALLBACK_DONE {",serial,"}")
}
