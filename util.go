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
	var id int32
	id = -1
	for k,val := range objects {
		if val == nil && k != 0 {
			id = k
		}
	}
	if id == -1 {
		id = int32(len(objects))
	}
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
	objects[id] = nil
}

func PrintObject(id int32) {
	fmt.Printf("Object id is %d\n", objects[id].Id())
}

func printError(f string, err error) {
	fmt.Println(f,"produced an error:",err)
}

func printEvent(name string, event string, args ...interface{}) {
	fmt.Printf("%s.%s { %v }\n",name,event,args)
}

func printRequest(name string, req string, args ...interface{}) {
	fmt.Printf(" -> %s.%s { %v }\n",name,req,args)
}

func delete_id_listener(c chan interface{}) {
	for e := range c {
		ev := e.(DisplayDelete_id)
		removeObject(int32(ev.Id))
	}
}

func error_listener(c chan interface{}) {
	for e := range c {
		ev := e.(DisplayError)
		fmt.Println("Error:", ev.Object_id.Id(), ev.Code, ev.Message)
	}
}
