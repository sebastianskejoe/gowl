
package gowl

import (
	"bytes"
)

type Shm struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shm, msg []byte)
}

//// Requests
func (s *Shm) Create_pool (id *Shm_pool, fd uintptr, size int32 ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())
	writeInteger(buf, fd)
	writeInteger(buf, size)

	sendmsg(s, 0, buf.Bytes())
}

//// Events
func (s *Shm) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

type ShmFormat struct {
	Format uint32
}

func (s *Shm) AddFormatListener(channel chan interface{}) {
	s.listeners[0] = append(s.listeners[0], channel)
}

func shm_format(s *Shm, msg []byte) {
	printEvent("format", msg)
	var data ShmFormat
	buf := bytes.NewBuffer(msg)

	format,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Format = format

	for _,channel := range s.listeners[0] {
		go func () { channel <- data }()
	}
}

func NewShm() (s *Shm) {
	s = new(Shm)
	s.listeners = make(map[int16][]chan interface{}, 0)

	s.events = append(s.events, shm_format)
	return
}

func (s *Shm) SetId(id int32) {
	s.id = id
}

func (s *Shm) Id() int32 {
	return s.id
}