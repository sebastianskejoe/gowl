package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Shm struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shm, msg message)
}

//// Requests
func (s *Shm) CreatePool (id *ShmPool, fd uintptr, size int32) {
	msg := newMessage(s, 0)
	appendObject(id)
	writeInteger(msg,id.Id())
	writeFd(msg,fd)
	writeInteger(msg,size)

	sendmsg(msg)
	printRequest("shm", s, "create_pool", "new id", id.Id(), fd, size)
}

//// Events
func (s *Shm) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

type ShmFormat struct {
	Format uint32
}

func (s *Shm) AddFormatListener(channel chan interface{}) {
	s.listeners[0] = append(s.listeners[0], channel)
	addListener(channel)
}

func shm_format(s *Shm, msg message) {
	var data ShmFormat

	format,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Format = format

	for _,channel := range s.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("shm", s, "format", format)
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