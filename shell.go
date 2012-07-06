package gowl
type Shell struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shell, msg []byte)
}

//// Requests
func (s *Shell) Get_shell_surface (id *Shell_surface, surface *Surface) {
	msg := newMessage(s, 0)
	appendObject(id)
	writeInteger(msg,id.Id())
	writeInteger(msg,surface.Id())

	sendmsg(msg)
	printRequest("shell", "get_shell_surface", id, surface)
}

//// Events
func (s *Shell) HandleEvent(opcode int16, msg []byte) {
	if s.events[opcode] != nil {
		s.events[opcode](s, msg)
	}
}

func NewShell() (s *Shell) {
	s = new(Shell)
	s.listeners = make(map[int16][]chan interface{}, 0)

	return
}

func (s *Shell) SetId(id int32) {
	s.id = id
}

func (s *Shell) Id() int32 {
	return s.id
}