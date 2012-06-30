package gowl

import (
	"bytes"
)

type Compositor struct {
	events map[int16]func()
	id int32
}

func (c *Compositor) HandleEvent(opcode int16, msg []byte) {
	return
}

func (c *Compositor) SetID(id int32) {
	c.id = id
}

func (c *Compositor) ID() int32 {
	return c.id
}

func (c *Compositor) CreateSurface() (s *Surface) {
	s = new(Surface)
	appendObject(s)

	buf := new(bytes.Buffer)
	writeInteger(buf, s.ID())
	sendmsg(c, 0, buf.Bytes())

	printRequest("CREATE_SURFACE", s.ID())
	return
}
