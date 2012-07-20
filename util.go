package gowl

import (
	"fmt"
)

var objects map[int32]Object
var freeIds []int32
var idchan chan int32

func init() {
	objects = make(map[int32]Object)
	objects[0] = nil

	freeIds = make([]int32, 0)
	idchan = make(chan int32)
	go pushIds(idchan)
}

func pushIds(c chan int32) {
	var id int32
	for {
		if len(freeIds) > 0 {
			id, freeIds = freeIds[0], freeIds[1:]
		} else {
			id = int32(len(objects))
		}
		c <- id
	}
}

func appendObject(obj Object) int32 {
	id := <-idchan
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
	freeIds = append(freeIds, id)
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
		ev := e.(DisplayDeleteId)
		removeObject(int32(ev.Id))
	}
}

func error_listener(c chan interface{}) {
	for e := range c {
		ev := e.(DisplayError)
		fmt.Println("Error:", ev.Object_id.Id(), ev.Code, ev.Message)
	}
}
