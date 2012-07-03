
package gowl

import (
	"bytes"
)

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
	serial uint32
	time uint32
	surface *Surface
	id int32
	x int32
	y int32
}

func (t *Touch) AddDownListener(channel chan interface{}) {
	t.addListener(0, channel)
}

func touch_down(t *Touch, msg []byte) {
	printEvent("down", msg)
	var data TouchDown
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.serial = serial

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.time = time

	surfaceid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surface = getObject(surfaceid).(*Surface)
	data.surface = surface

	id,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.id = id

	x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.x = x

	y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.y = y

	for _,channel := range t.listeners[0] {
		channel <- data
	}
}

type TouchUp struct {
	serial uint32
	time uint32
	id int32
}

func (t *Touch) AddUpListener(channel chan interface{}) {
	t.addListener(1, channel)
}

func touch_up(t *Touch, msg []byte) {
	printEvent("up", msg)
	var data TouchUp
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.serial = serial

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.time = time

	id,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.id = id

	for _,channel := range t.listeners[1] {
		channel <- data
	}
}

type TouchMotion struct {
	time uint32
	id int32
	x int32
	y int32
}

func (t *Touch) AddMotionListener(channel chan interface{}) {
	t.addListener(2, channel)
}

func touch_motion(t *Touch, msg []byte) {
	printEvent("motion", msg)
	var data TouchMotion
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.time = time

	id,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.id = id

	x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.x = x

	y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.y = y

	for _,channel := range t.listeners[2] {
		channel <- data
	}
}

type TouchFrame struct {
}

func (t *Touch) AddFrameListener(channel chan interface{}) {
	t.addListener(3, channel)
}

func touch_frame(t *Touch, msg []byte) {
	printEvent("frame", msg)
	var data TouchFrame

	for _,channel := range t.listeners[3] {
		channel <- data
	}
}

type TouchCancel struct {
}

func (t *Touch) AddCancelListener(channel chan interface{}) {
	t.addListener(4, channel)
}

func touch_cancel(t *Touch, msg []byte) {
	printEvent("cancel", msg)
	var data TouchCancel

	for _,channel := range t.listeners[4] {
		channel <- data
	}
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