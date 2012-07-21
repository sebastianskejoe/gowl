package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Touch struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (t *Touch, msg []byte)
}

//// Requests
//// Events
func (t *Touch) HandleEvent(opcode int16, msg []byte) {
	if t.events[opcode] != nil {
		t.events[opcode](t, msg)
	}
}

type TouchDown struct {
	Serial uint32
	Time uint32
	Surface *Surface
	Id int32
	X int32
	Y int32
}

func (t *Touch) AddDownListener(channel chan interface{}) {
	t.listeners[0] = append(t.listeners[0], channel)
}

func touch_down(t *Touch, msg []byte) {
	var data TouchDown
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	surfaceid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surfaceobj := getObject(surfaceid)
	if surfaceobj == nil {
		return
	}
	surface = surfaceobj.(*Surface)
	data.Surface = surface

	id,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Id = id

	x,err := readFixed(buf)
	if err != nil {
		// XXX Error handling
	}
	data.X = x

	y,err := readFixed(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Y = y

	for _,channel := range t.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("touch", "down", serial, time, surface, id, x, y)
}

type TouchUp struct {
	Serial uint32
	Time uint32
	Id int32
}

func (t *Touch) AddUpListener(channel chan interface{}) {
	t.listeners[1] = append(t.listeners[1], channel)
}

func touch_up(t *Touch, msg []byte) {
	var data TouchUp
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	id,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Id = id

	for _,channel := range t.listeners[1] {
		go func () { channel <- data }()
	}
	printEvent("touch", "up", serial, time, id)
}

type TouchMotion struct {
	Time uint32
	Id int32
	X int32
	Y int32
}

func (t *Touch) AddMotionListener(channel chan interface{}) {
	t.listeners[2] = append(t.listeners[2], channel)
}

func touch_motion(t *Touch, msg []byte) {
	var data TouchMotion
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	id,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Id = id

	x,err := readFixed(buf)
	if err != nil {
		// XXX Error handling
	}
	data.X = x

	y,err := readFixed(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Y = y

	for _,channel := range t.listeners[2] {
		go func () { channel <- data }()
	}
	printEvent("touch", "motion", time, id, x, y)
}

type TouchFrame struct {
}

func (t *Touch) AddFrameListener(channel chan interface{}) {
	t.listeners[3] = append(t.listeners[3], channel)
}

func touch_frame(t *Touch, msg []byte) {
	var data TouchFrame

	for _,channel := range t.listeners[3] {
		go func () { channel <- data }()
	}
	printEvent("touch", "frame", )
}

type TouchCancel struct {
}

func (t *Touch) AddCancelListener(channel chan interface{}) {
	t.listeners[4] = append(t.listeners[4], channel)
}

func touch_cancel(t *Touch, msg []byte) {
	var data TouchCancel

	for _,channel := range t.listeners[4] {
		go func () { channel <- data }()
	}
	printEvent("touch", "cancel", )
}

func NewTouch() (t *Touch) {
	t = new(Touch)
	t.listeners = make(map[int16][]chan interface{}, 0)

	t.events = append(t.events, touch_down)
	t.events = append(t.events, touch_up)
	t.events = append(t.events, touch_motion)
	t.events = append(t.events, touch_frame)
	t.events = append(t.events, touch_cancel)
	return
}

func (t *Touch) SetId(id int32) {
	t.id = id
}

func (t *Touch) Id() int32 {
	return t.id
}