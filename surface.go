package gowl

type Surface struct {
	events map[int16]func()
	id int32
}

func (s *Surface) HandleEvent(opcode int16, msg []byte) {
	return
}

func (s *Surface) SetID(id int32) {
	s.id = id
}

func (s *Surface) ID() int32 {
	return s.id
}
