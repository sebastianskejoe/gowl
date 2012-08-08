package gowl

type Compositor struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (c *Compositor, msg message)
}

//// Requests
func (c *Compositor) CreateSurface (id *Surface) {
	msg := newMessage(c, 0)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("compositor", c, "create_surface", "new id", id.Id())
}

func (c *Compositor) CreateRegion (id *Region) {
	msg := newMessage(c, 1)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("compositor", c, "create_region", "new id", id.Id())
}

//// Events
func (c *Compositor) HandleEvent(msg message) {
	if c.events[msg.opcode] != nil {
		c.events[msg.opcode](c, msg)
	}
}

func NewCompositor() (c *Compositor) {
	c = new(Compositor)
	c.listeners = make(map[int16][]chan interface{}, 0)

	return
}

func (c *Compositor) SetId(id int32) {
	c.id = id
}

func (c *Compositor) Id() int32 {
	return c.id
}