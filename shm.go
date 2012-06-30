package gowl

import (
	"fmt"
	"bytes"
)

type Shm struct {
	events map[int16]func()
	id int32
}

func (s *Shm) HandleEvent(opcode int16, msg []byte) {
	if opcode != 0 {
		fmt.Println("Unknown shm event opcode",opcode)
		return
	}
	shm_format(msg)
	return
}

func (s *Shm) SetID(id int32) {
	s.id = id
}

func (s *Shm) ID() int32 {
	return s.id
}

func shm_format(msg []byte) {
	format, err := readUint32(bytes.NewBuffer(msg))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("SHM_FORMAT {",format,"}")
}
