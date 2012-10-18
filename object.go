package gowl

type WlObject struct {
	id int32
	listeners map[int16][]chan interface{}//func(...interface{})
}

func (obj *WlObject) SetId(id int32) {
	obj.id = id
}

func (obj *WlObject) Id() int32 {
	return obj.id
}

func (obj *WlObject) addListener(opcode int16, c chan interface{}) {
	obj.listeners[opcode] = append(obj.listeners[opcode], c)
}

type Object interface {
	HandleEvent(msg message)
	SetId(id int32)
	Id() int32
    Name() string
}

