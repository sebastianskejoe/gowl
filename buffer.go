
package gowl

import (
	"bytes"
)

type Buffer struct {
	*WlObject
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
	b.addListener(0, channel)
}

func buffer_release(b *Buffer, msg []byte) {
	var data BufferRelease

	for _,channel := range b.listeners[0] {
		channel <- data
	}
}

func NewBuffer() (b *Buffer) {
	b = new(Buffer)

	b.events = append(b.events, buffer_release)
	return
}