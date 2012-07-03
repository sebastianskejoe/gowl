
package gowl

import (
	"bytes"
)

type Compositor struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (c *Compositor, msg []byte)
}

//// Requests
func (c *Compositor) Create_surface (id *Surface ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(c, 0, buf.Bytes())
}

func (c *Compositor) Create_region (id *Region ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(c, 1, buf.Bytes())
}

//// Events
func (c *Compositor) HandleEvent(opcode int16, msg []byte) {
	if c.events[opcode] != nil {
		c.events[opcode](c, msg)
	}
}

func NewCompositor() (c *Compositor) {
	c = new(Compositor)
	c.listeners = make(map[int16][]chan interface{}, 0)

	return
}