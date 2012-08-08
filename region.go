package gowl

type Region struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (r *Region, msg message)
}

//// Requests
func (r *Region) Destroy () {
	msg := newMessage(r, 0)

	sendmsg(msg)
	printRequest("region", r, "destroy")
}

func (r *Region) Add (x int32, y int32, width int32, height int32) {
	msg := newMessage(r, 1)
	writeInteger(msg,x)
	writeInteger(msg,y)
	writeInteger(msg,width)
	writeInteger(msg,height)

	sendmsg(msg)
	printRequest("region", r, "add", x, y, width, height)
}

func (r *Region) Subtract (x int32, y int32, width int32, height int32) {
	msg := newMessage(r, 2)
	writeInteger(msg,x)
	writeInteger(msg,y)
	writeInteger(msg,width)
	writeInteger(msg,height)

	sendmsg(msg)
	printRequest("region", r, "subtract", x, y, width, height)
}

//// Events
func (r *Region) HandleEvent(msg message) {
	if r.events[msg.opcode] != nil {
		r.events[msg.opcode](r, msg)
	}
}

func NewRegion() (r *Region) {
	r = new(Region)
	r.listeners = make(map[int16][]chan interface{}, 0)

	return
}

func (r *Region) SetId(id int32) {
	r.id = id
}

func (r *Region) Id() int32 {
	return r.id
}