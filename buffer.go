
package gowl

import (
	"bytes"
)

type Buffer struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (b *Buffer, msg []byte)
}

//// Requests
func (b *Buffer) Destroy ( ) {
	buf := new(bytes.Buffer)

	sendmsg(b, 0, buf.Bytes())
}

//// Events
func (b *Buffer) HandleEvent(opcode int16, msg []byte) {
	if b.events[opcode] != nil {
		b.events[opcode](b, msg)
	}
}

type BufferRelease struct {
}

func (b *Buffer) AddReleaseListener(channel chan interface{}) {
	b.listeners[0] = append(b.listeners[0], channel)
}

func buffer_release(b *Buffer, msg []byte) {
	printEvent("release", msg)
	var data BufferRelease

	for _,channel := range b.listeners[0] {
		channel <- data
	}
}

func NewBuffer() (b *Buffer) {
	b = new(Buffer)
	b.listeners = make(map[int16][]chan interface{}, 0)

	b.events = append(b.events, buffer_release)
	return
}

func (b *Buffer) SetId(id int32) {
	b.id = id
}

func (b *Buffer) Id() int32 {
	return b.id
}