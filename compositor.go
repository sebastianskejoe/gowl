
package gowl

import (
	"bytes"
)

type Compositor struct {
	*WlObject
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

	return
}