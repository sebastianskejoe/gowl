package gowl
type Compositor struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (c *Compositor, msg []byte)
}

//// Requests
func (c *Compositor) Create_surface (id *Surface) {
	msg := newMessage(c, 0)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("compositor", "create_surface", id)
}

func (c *Compositor) Create_region (id *Region) {
	msg := newMessage(c, 1)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("compositor", "create_region", id)
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

func (c *Compositor) SetId(id int32) {
	c.id = id
}

func (c *Compositor) Id() int32 {
	return c.id
}