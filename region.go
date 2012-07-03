
package gowl

import (
	"bytes"
)

type Region struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (r *Region, msg []byte)
}

//// Requests
func (r *Region) Destroy ( ) {
	buf := new(bytes.Buffer)

	sendmsg(r, 0, buf.Bytes())
}

func (r *Region) Add (x int32, y int32, width int32, height int32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, x)
	writeInteger(buf, y)
	writeInteger(buf, width)
	writeInteger(buf, height)

	sendmsg(r, 1, buf.Bytes())
}

func (r *Region) Subtract (x int32, y int32, width int32, height int32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, x)
	writeInteger(buf, y)
	writeInteger(buf, width)
	writeInteger(buf, height)

	sendmsg(r, 2, buf.Bytes())
}

//// Events
func (r *Region) HandleEvent(opcode int16, msg []byte) {
	if r.events[opcode] != nil {
		r.events[opcode](r, msg)
	}
}

func NewRegion() (r *Region) {
	r = new(Region)
	r.listeners = make(map[int16][]chan interface{}, 0)

	return
}