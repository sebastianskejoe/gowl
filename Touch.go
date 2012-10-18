package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Touch struct {
	id int32
    downListeners []chan TouchDown
	upListeners []chan TouchUp
	motionListeners []chan TouchMotion
	frameListeners []chan TouchFrame
	cancelListeners []chan TouchCancel
	events []func(t *Touch, msg message) error
    name string
}

func NewTouch() (t *Touch) {
	t = new(Touch)
    t.name = "Touch"
    t.downListeners = make([]chan TouchDown, 0)
	t.upListeners = make([]chan TouchUp, 0)
	t.motionListeners = make([]chan TouchMotion, 0)
	t.frameListeners = make([]chan TouchFrame, 0)
	t.cancelListeners = make([]chan TouchCancel, 0)

    t.events = append(t.events, touchDown)
	t.events = append(t.events, touchUp)
	t.events = append(t.events, touchMotion)
	t.events = append(t.events, touchFrame)
	t.events = append(t.events, touchCancel)
	return
}

func (t *Touch) HandleEvent(msg message) {
	if t.events[msg.opcode] != nil {
		t.events[msg.opcode](t, msg)
	}
}

func (t *Touch) SetId(id int32) {
	t.id = id
}

func (t *Touch) Id() int32 {
	return t.id
}

func (t *Touch) Name() string {
    return t.name
}

////
//// REQUESTS
////

////
//// EVENTS
////

type TouchDown struct {
    Serial uint32
	Time uint32
	Surface *Surface
	Id int32
	X int32
	Y int32
}

func (t *Touch) AddDownListener(channel chan TouchDown) {
    t.downListeners = append(t.downListeners, channel)
}

func touchDown(t *Touch, msg message) (err error) {
    var data TouchDown

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Read time
    time,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Time = time

    // Read surface
    surface,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    surfaceObj := getObject(surface)
    data.Surface = surfaceObj.(*Surface)

    // Read id
    id,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Id = id

    // Read x
    x,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.X = x

    // Read y
    y,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.Y = y

    // Dispatch events
    for _,channel := range t.downListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type TouchUp struct {
    Serial uint32
	Time uint32
	Id int32
}

func (t *Touch) AddUpListener(channel chan TouchUp) {
    t.upListeners = append(t.upListeners, channel)
}

func touchUp(t *Touch, msg message) (err error) {
    var data TouchUp

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Read time
    time,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Time = time

    // Read id
    id,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Id = id

    // Dispatch events
    for _,channel := range t.upListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type TouchMotion struct {
    Time uint32
	Id int32
	X int32
	Y int32
}

func (t *Touch) AddMotionListener(channel chan TouchMotion) {
    t.motionListeners = append(t.motionListeners, channel)
}

func touchMotion(t *Touch, msg message) (err error) {
    var data TouchMotion

    // Read time
    time,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Time = time

    // Read id
    id,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    data.Id = id

    // Read x
    x,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.X = x

    // Read y
    y,err := readFixed(msg.buf)
    if err != nil {
        return
    }
    data.Y = y

    // Dispatch events
    for _,channel := range t.motionListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type TouchFrame struct {
    
}

func (t *Touch) AddFrameListener(channel chan TouchFrame) {
    t.frameListeners = append(t.frameListeners, channel)
}

func touchFrame(t *Touch, msg message) (err error) {
    var data TouchFrame


    // Dispatch events
    for _,channel := range t.frameListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type TouchCancel struct {
    
}

func (t *Touch) AddCancelListener(channel chan TouchCancel) {
    t.cancelListeners = append(t.cancelListeners, channel)
}

func touchCancel(t *Touch, msg message) (err error) {
    var data TouchCancel


    // Dispatch events
    for _,channel := range t.cancelListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
