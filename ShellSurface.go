package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type ShellSurface struct {
	id int32
    pingListeners []chan ShellSurfacePing
	configureListeners []chan ShellSurfaceConfigure
	popupDoneListeners []chan ShellSurfacePopupDone
	events []func(s *ShellSurface, msg message) error
    name string
}

func NewShellSurface() (s *ShellSurface) {
	s = new(ShellSurface)
    s.name = "ShellSurface"
    s.pingListeners = make([]chan ShellSurfacePing, 0)
	s.configureListeners = make([]chan ShellSurfaceConfigure, 0)
	s.popupDoneListeners = make([]chan ShellSurfacePopupDone, 0)

    s.events = append(s.events, shellSurfacePing)
	s.events = append(s.events, shellSurfaceConfigure)
	s.events = append(s.events, shellSurfacePopupDone)
	return
}

func (s *ShellSurface) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

func (s *ShellSurface) SetId(id int32) {
	s.id = id
}

func (s *ShellSurface) Id() int32 {
	return s.id
}

func (s *ShellSurface) Name() string {
    return s.name
}

////
//// REQUESTS
////

func (s *ShellSurface) Pong(serial uint32) {
    sendrequest(s, "wl_shell_surface_pong", serial)
}

func (s *ShellSurface) Move(seat *Seat, serial uint32) {
    sendrequest(s, "wl_shell_surface_move", seat, serial)
}

func (s *ShellSurface) Resize(seat *Seat, serial uint32, edges uint32) {
    sendrequest(s, "wl_shell_surface_resize", seat, serial, edges)
}

func (s *ShellSurface) SetToplevel() {
    sendrequest(s, "wl_shell_surface_set_toplevel", )
}

func (s *ShellSurface) SetTransient(parent *Surface, x int32, y int32, flags uint32) {
    sendrequest(s, "wl_shell_surface_set_transient", parent, x, y, flags)
}

func (s *ShellSurface) SetFullscreen(method uint32, framerate uint32, output *Output) {
    sendrequest(s, "wl_shell_surface_set_fullscreen", method, framerate, output)
}

func (s *ShellSurface) SetPopup(seat *Seat, serial uint32, parent *Surface, x int32, y int32, flags uint32) {
    sendrequest(s, "wl_shell_surface_set_popup", seat, serial, parent, x, y, flags)
}

func (s *ShellSurface) SetMaximized(output *Output) {
    sendrequest(s, "wl_shell_surface_set_maximized", output)
}

func (s *ShellSurface) SetTitle(title string) {
    sendrequest(s, "wl_shell_surface_set_title", title)
}

func (s *ShellSurface) SetClass(class_ string) {
    sendrequest(s, "wl_shell_surface_set_class", class_)
}

////
//// EVENTS
////

type ShellSurfacePing struct {
    Serial uint32
}

func (s *ShellSurface) AddPingListener(channel chan ShellSurfacePing) {
    s.pingListeners = append(s.pingListeners, channel)
}

func shellSurfacePing(s *ShellSurface, msg message) (err error) {
    var data ShellSurfacePing

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Dispatch events
    for _,channel := range s.pingListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type ShellSurfaceConfigure struct {
    Edges uint32
	Width int32
	Height int32
}

func (s *ShellSurface) AddConfigureListener(channel chan ShellSurfaceConfigure) {
    s.configureListeners = append(s.configureListeners, channel)
}

func shellSurfaceConfigure(s *ShellSurface, msg message) (err error) {
    var data ShellSurfaceConfigure

    // Read edges
    edges,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Edges = edges

    // Read width
    width,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Width = width

    // Read height
    height,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Height = height

    // Dispatch events
    for _,channel := range s.configureListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type ShellSurfacePopupDone struct {
    
}

func (s *ShellSurface) AddPopupDoneListener(channel chan ShellSurfacePopupDone) {
    s.popupDoneListeners = append(s.popupDoneListeners, channel)
}

func shellSurfacePopupDone(s *ShellSurface, msg message) (err error) {
    var data ShellSurfacePopupDone


    // Dispatch events
    for _,channel := range s.popupDoneListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
