package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Shell struct {
	id int32
    
	events []func(s *Shell, msg message) error
    name string
}

func NewShell() (s *Shell) {
	s = new(Shell)
    s.name = "Shell"
    

    
	return
}

func (s *Shell) HandleEvent(msg message) {
	if s.events[msg.opcode] != nil {
		s.events[msg.opcode](s, msg)
	}
}

func (s *Shell) SetId(id int32) {
	s.id = id
}

func (s *Shell) Id() int32 {
	return s.id
}

func (s *Shell) Name() string {
    return s.name
}

////
//// REQUESTS
////

func (s *Shell) GetShellSurface(id *ShellSurface, surface *Surface) {
    sendrequest(s, "wl_shell_get_shell_surface", id, surface)
}

////
//// EVENTS
////
