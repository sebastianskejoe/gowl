package gowl

import (
	"fmt"
	"io/ioutil"
	"os"
)

var objects map[int32]Object
var freeIds []int32
var idchan chan int32
var nextid int32

func init() {
	objects = make(map[int32]Object)
	objects[0] = nil

	freeIds = make([]int32, 0)
	idchan = make(chan int32)
	nextid = 1
	go pushIds(idchan)
}

func pushIds(c chan int32) {
	var id int32
	for {
		if len(freeIds) > 0 {
			id, freeIds = freeIds[0], freeIds[1:]
		} else {
			id = nextid
			nextid++
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
		fmt.Println("Error:", ev.ObjectId.Id(), ev.Code, ev.Message)
	}
}

func (d *Display) Iterate() error {
	for {
		id, opcode, _, msg, remain, err := getmsg()
		if err != nil {
			return err
		}
		obj := getObject(id)
		if obj != nil {
			obj.HandleEvent(opcode, msg)
		}

		if remain == 0 {
			break
		}
	}
	return nil
}

func (d *Display) Connect() error {
	err := connect_to_socket()
	if err != nil {
		return err
	}
	appendObject(d)

	delchan := make(chan interface{})
	d.AddDeleteIdListener(delchan)
	errchan := make(chan interface{})
	d.AddErrorListener(errchan)
	go delete_id_listener(delchan)
	go error_listener(errchan)
	return nil
}

func CreateTmp(size int64) (uintptr) {
	tmp,err := ioutil.TempFile(os.Getenv("XDG_RUNTIME_DIR"), "gowl")
	if err != nil {
		fmt.Println(err)
	}
	tmp.Truncate(size)
	os.Remove(tmp.Name())
	return tmp.Fd()
}
