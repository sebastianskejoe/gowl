package gowl

type Shell struct {
	events map[int16]func()
	id int32
}

func (s *Shell) HandleEvent(opcode int16, msg []byte) {
	return
}

func (s *Shell) SetID(id int32) {
	s.id = id
}

func (s *Shell) ID() int32 {
	return s.id
}
