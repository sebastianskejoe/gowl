package gowl

type Shell struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (s *Shell, msg message)
}

//// Requests
func (s *Shell) GetShellSurface (id *ShellSurface, surface *Surface) {
	msg := newMessage(s, 0)
	appendObject(id)
	writeInteger(msg,id.Id())
	writeInteger(msg,surface.Id())

	sendmsg(msg)
	printRequest("shell", s, "get_shell_surface", "new id", id.Id(), surface.Id())
}

//// Events
func (s *Shell) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
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