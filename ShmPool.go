package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type ShmPool struct {
	id int32
    
	events []func(s *ShmPool, msg message) error
    name string
}

func NewShmPool() (s *ShmPool) {
	s = new(ShmPool)
    s.name = "ShmPool"
    

    
	return
}

func (s *ShmPool) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

func (s *ShmPool) SetId(id int32) {
	s.id = id
}

func (s *ShmPool) Id() int32 {
	return s.id
}

func (s *ShmPool) Name() string {
    return s.name
}

////
//// REQUESTS
////

func (s *ShmPool) CreateBuffer(id *Buffer, offset int32, width int32, height int32, stride int32, format uint32) {
    sendrequest(s, "wl_shm_pool_create_buffer", id, offset, width, height, stride, format)
}

func (s *ShmPool) Destroy() {
    sendrequest(s, "wl_shm_pool_destroy", )
}

func (s *ShmPool) Resize(size int32) {
    sendrequest(s, "wl_shm_pool_resize", size)
}

////
//// EVENTS
////
