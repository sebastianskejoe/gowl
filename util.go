package gowl

import (
	"fmt"
)

var __id int32

type Object interface {
	HandleEvent(opcode int16, msg []byte)
	SetID(id int32)
	ID() int32
}

var objects map[int32]Object

func init() {
	__id = 0

	objects = make(map[int32]Object)
	objects[0] = nil
}

func appendObject(obj Object) int32 {
	id := int32(len(objects))
	objects[id] = obj
	obj.SetID(id)
	return id
}

func PrintObject(id int32) {
	fmt.Printf("%d\n", objects[id].ID())
}
