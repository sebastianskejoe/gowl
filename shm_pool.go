package gowl

type ShmPool struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *ShmPool, msg message)
}

//// Requests
func (s *ShmPool) CreateBuffer (id *Buffer, offset int32, width int32, height int32, stride int32, format uint32) {
	msg := newMessage(s, 0)
	appendObject(id)
	writeInteger(msg,id.Id())
	writeInteger(msg,offset)
	writeInteger(msg,width)
	writeInteger(msg,height)
	writeInteger(msg,stride)
	writeInteger(msg,format)

	sendmsg(msg)
	printRequest("shm_pool", s, "create_buffer", "new id", id.Id(), offset, width, height, stride, format)
}

func (s *ShmPool) Destroy () {
	msg := newMessage(s, 1)

	sendmsg(msg)
	printRequest("shm_pool", s, "destroy")
}

func (s *ShmPool) Resize (size int32) {
	msg := newMessage(s, 2)
	writeInteger(msg,size)

	sendmsg(msg)
	printRequest("shm_pool", s, "resize", size)
}

//// Events
func (s *ShmPool) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

func NewShmPool() (s *ShmPool) {
	s = new(ShmPool)
	s.listeners = make(map[int16][]chan interface{}, 0)

	return
}

func (s *ShmPool) SetId(id int32) {
	s.id = id
}

func (s *ShmPool) Id() int32 {
	return s.id
}