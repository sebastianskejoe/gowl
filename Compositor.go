package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Compositor struct {
	id int32
    
	events []func(c *Compositor, msg message) error
    name string
}

func NewCompositor() (c *Compositor) {
	c = new(Compositor)
    c.name = "Compositor"
    

    
	return
}

func (c *Compositor) HandleEvent(msg message) {
	if c.events[msg.opcode] != nil {
		c.events[msg.opcode](c, msg)
	}
}

func (c *Compositor) SetId(id int32) {
	c.id = id
}

func (c *Compositor) Id() int32 {
	return c.id
}

func (c *Compositor) Name() string {
    return c.name
}

////
//// REQUESTS
////

func (c *Compositor) CreateSurface(id *Surface) {
    sendrequest(c, "wl_compositor_create_surface", id)
}

func (c *Compositor) CreateRegion(id *Region) {
    sendrequest(c, "wl_compositor_create_region", id)
}

////
//// EVENTS
////
