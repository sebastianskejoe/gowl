package gowl

import (
	"fmt"
)

var __id int32


var objects map[int32]Object

func init() {
	objects = make(map[int32]Object)
	objects[0] = nil
}

func appendObject(obj Object) int32 {
	id := int32(len(objects))
	objects[id] = obj
	obj.SetId(id)
	return id
}

func setObject(id int32, obj Object) {
	objects[id] = obj
}

func getObject(id int32) Object {
	return objects[id]
}

func removeObject(id int32) {
	delete(objects, id)
}

func PrintObject(id int32) {
	fmt.Printf("%d\n", objects[id].Id())
}

func printError(f string, err error) {
	fmt.Println(f,"produced an error:",err)
}

func printRequest(name string, args ...interface{}) {
	fmt.Println("->",name,"{",args,"}")
}
