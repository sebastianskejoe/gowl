
package gowl

import (
	"bytes"
)

type Shm_pool struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shm_pool, msg []byte)
}

//// Requests
func (s *Shm_pool) Create_buffer (id *Buffer, offset int32, width int32, height int32, stride int32, format uint32 ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())
	writeInteger(buf, offset)
	writeInteger(buf, width)
	writeInteger(buf, height)
	writeInteger(buf, stride)
	writeInteger(buf, format)

	sendmsg(s, 0, buf.Bytes())
}

func (s *Shm_pool) Destroy ( ) {
	buf := new(bytes.Buffer)

	sendmsg(s, 1, buf.Bytes())
}

func (s *Shm_pool) Resize (size int32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, size)

	sendmsg(s, 2, buf.Bytes())
}

//// Events
func (s *Shm_pool) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

func NewShm_pool() (s *Shm_pool) {
	s = new(Shm_pool)
	s.listeners = make(map[int16][]chan interface{}, 0)

	return
}

func (s *Shm_pool) SetId(id int32) {
	s.id = id
}

func (s *Shm_pool) Id() int32 {
	return s.id
}