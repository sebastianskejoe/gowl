package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Region struct {
	id int32
    
	events []func(r *Region, msg message) error
    name string
}

func NewRegion() (r *Region) {
	r = new(Region)
    r.name = "Region"
    

    
	return
}

func (r *Region) HandleEvent(msg message) {
	if r.events[msg.opcode] != nil {
		r.events[msg.opcode](r, msg)
	}
}

func (r *Region) SetId(id int32) {
	r.id = id
}

func (r *Region) Id() int32 {
	return r.id
}

func (r *Region) Name() string {
    return r.name
}

////
//// REQUESTS
////

func (r *Region) Destroy() {
    sendrequest(r, "wl_region_destroy", )
}

func (r *Region) Add(x int32, y int32, width int32, height int32) {
    sendrequest(r, "wl_region_add", x, y, width, height)
}

func (r *Region) Subtract(x int32, y int32, width int32, height int32) {
    sendrequest(r, "wl_region_subtract", x, y, width, height)
}

////
//// EVENTS
////
