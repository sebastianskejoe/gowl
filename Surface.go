package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Surface struct {
	id int32
    enterListeners []chan SurfaceEnter
	leaveListeners []chan SurfaceLeave
	events []func(s *Surface, msg message) error
    name string
}

func NewSurface() (s *Surface) {
	s = new(Surface)
    s.name = "Surface"
    s.enterListeners = make([]chan SurfaceEnter, 0)
	s.leaveListeners = make([]chan SurfaceLeave, 0)

    s.events = append(s.events, surfaceEnter)
	s.events = append(s.events, surfaceLeave)
	return
}

func (s *Surface) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

func (s *Surface) SetId(id int32) {
	s.id = id
}

func (s *Surface) Id() int32 {
	return s.id
}

func (s *Surface) Name() string {
    return s.name
}

////
//// REQUESTS
////

func (s *Surface) Destroy() {
    sendrequest(s, "wl_surface_destroy", )
}

func (s *Surface) Attach(buffer *Buffer, x int32, y int32) {
    sendrequest(s, "wl_surface_attach", buffer, x, y)
}

func (s *Surface) Damage(x int32, y int32, width int32, height int32) {
    sendrequest(s, "wl_surface_damage", x, y, width, height)
}

func (s *Surface) Frame(callback *Callback) {
    sendrequest(s, "wl_surface_frame", callback)
}

func (s *Surface) SetOpaqueRegion(region *Region) {
    sendrequest(s, "wl_surface_set_opaque_region", region)
}

func (s *Surface) SetInputRegion(region *Region) {
    sendrequest(s, "wl_surface_set_input_region", region)
}

func (s *Surface) Commit() {
    sendrequest(s, "wl_surface_commit", )
}

////
//// EVENTS
////

type SurfaceEnter struct {
    Output *Output
}

func (s *Surface) AddEnterListener(channel chan SurfaceEnter) {
    s.enterListeners = append(s.enterListeners, channel)
}

func surfaceEnter(s *Surface, msg message) (err error) {
    var data SurfaceEnter

    // Read output
    output,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    outputObj := getObject(output)
    data.Output = outputObj.(*Output)

    // Dispatch events
    for _,channel := range s.enterListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type SurfaceLeave struct {
    Output *Output
}

func (s *Surface) AddLeaveListener(channel chan SurfaceLeave) {
    s.leaveListeners = append(s.leaveListeners, channel)
}

func surfaceLeave(s *Surface, msg message) (err error) {
    var data SurfaceLeave

    // Read output
    output,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    outputObj := getObject(output)
    data.Output = outputObj.(*Output)

    // Dispatch events
    for _,channel := range s.leaveListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
