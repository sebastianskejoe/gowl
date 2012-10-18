package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Shm struct {
	id int32
    formatListeners []chan ShmFormat
	events []func(s *Shm, msg message) error
    name string
}

func NewShm() (s *Shm) {
	s = new(Shm)
    s.name = "Shm"
    s.formatListeners = make([]chan ShmFormat, 0)

    s.events = append(s.events, shmFormat)
	return
}

func (s *Shm) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

func (s *Shm) SetId(id int32) {
	s.id = id
}

func (s *Shm) Id() int32 {
	return s.id
}

func (s *Shm) Name() string {
    return s.name
}

////
//// REQUESTS
////

func (s *Shm) CreatePool(id *ShmPool, fd uintptr, size int32) {
    sendrequest(s, "wl_shm_create_pool", id, fd, size)
}

////
//// EVENTS
////

type ShmFormat struct {
    Format uint32
}

func (s *Shm) AddFormatListener(channel chan ShmFormat) {
    s.formatListeners = append(s.formatListeners, channel)
}

func shmFormat(s *Shm, msg message) (err error) {
    var data ShmFormat

    // Read format
    format,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Format = format

    // Dispatch events
    for _,channel := range s.formatListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
