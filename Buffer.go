package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Buffer struct {
	id int32
    releaseListeners []chan BufferRelease
	events []func(b *Buffer, msg message) error
    name string
}

func NewBuffer() (b *Buffer) {
	b = new(Buffer)
    b.name = "Buffer"
    b.releaseListeners = make([]chan BufferRelease, 0)

    b.events = append(b.events, bufferRelease)
	return
}

func (b *Buffer) HandleEvent(msg message) {
	if b.events[msg.opcode] != nil {
		b.events[msg.opcode](b, msg)
	}
}

func (b *Buffer) SetId(id int32) {
	b.id = id
}

func (b *Buffer) Id() int32 {
	return b.id
}

func (b *Buffer) Name() string {
    return b.name
}

////
//// REQUESTS
////

func (b *Buffer) Destroy() {
    sendrequest(b, "wl_buffer_destroy", )
}

////
//// EVENTS
////

type BufferRelease struct {
    
}

func (b *Buffer) AddReleaseListener(channel chan BufferRelease) {
    b.releaseListeners = append(b.releaseListeners, channel)
}

func bufferRelease(b *Buffer, msg message) (err error) {
    var data BufferRelease


    // Dispatch events
    for _,channel := range b.releaseListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
