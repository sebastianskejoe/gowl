package gowl
type Shm_pool struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shm_pool, msg []byte)
}

//// Requests
func (s *Shm_pool) Create_buffer (id *Buffer, offset int32, width int32, height int32, stride int32, format uint32) {
	msg := newMessage(s, 0)
	appendObject(id)
	writeInteger(msg,id.Id())
	writeInteger(msg,offset)
	writeInteger(msg,width)
	writeInteger(msg,height)
	writeInteger(msg,stride)
	writeInteger(msg,format)

	sendmsg(msg)
	printRequest("shm_pool", "create_buffer", id, offset, width, height, stride, format)
}

func (s *Shm_pool) Destroy () {
	msg := newMessage(s, 1)

	sendmsg(msg)
	printRequest("shm_pool", "destroy", )
}

func (s *Shm_pool) Resize (size int32) {
	msg := newMessage(s, 2)
	writeInteger(msg,size)

	sendmsg(msg)
	printRequest("shm_pool", "resize", size)
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