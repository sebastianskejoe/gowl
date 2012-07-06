package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type Data_device struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Data_device, msg []byte)
}

//// Requests
func (d *Data_device) Start_drag (source *Data_source, origin *Surface, icon *Surface, serial uint32) {
	msg := newMessage(d, 0)
	writeInteger(msg,source.Id())
	writeInteger(msg,origin.Id())
	writeInteger(msg,icon.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("data_device", "start_drag", source, origin, icon, serial)
}

func (d *Data_device) Set_selection (source *Data_source, serial uint32) {
	msg := newMessage(d, 1)
	writeInteger(msg,source.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("data_device", "set_selection", source, serial)
}

//// Events
func (d *Data_device) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type Data_deviceData_offer struct {
	Id *Data_offer
}

func (d *Data_device) AddData_offerListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
}

func data_device_data_offer(d *Data_device, msg []byte) {
	var data Data_deviceData_offer
	buf := bytes.NewBuffer(msg)

	idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(Data_offer)
	setObject(idid, id)
	data.Id = id

	for _,channel := range d.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "data_offer", id)
}

type Data_deviceEnter struct {
	Serial uint32
	Surface *Surface
	X int32
	Y int32
	Id *Data_offer
}

func (d *Data_device) AddEnterListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
}

func data_device_enter(d *Data_device, msg []byte) {
	var data Data_deviceEnter
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

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

	x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.X = x

	y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Y = y

	idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(Data_offer)
	idobj := getObject(idid)
	if idobj == nil {
		return
	}
	id = idobj.(*Data_offer)
	data.Id = id

	for _,channel := range d.listeners[1] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "enter", serial, surface, x, y, id)
}

type Data_deviceLeave struct {
}

func (d *Data_device) AddLeaveListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
}

func data_device_leave(d *Data_device, msg []byte) {
	var data Data_deviceLeave

	for _,channel := range d.listeners[2] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "leave", )
}

type Data_deviceMotion struct {
	Time uint32
	X int32
	Y int32
}

func (d *Data_device) AddMotionListener(channel chan interface{}) {
	d.listeners[3] = append(d.listeners[3], channel)
}

func data_device_motion(d *Data_device, msg []byte) {
	var data Data_deviceMotion
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	x,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.X = x

	y,err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.Y = y

	for _,channel := range d.listeners[3] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "motion", time, x, y)
}

type Data_deviceDrop struct {
}

func (d *Data_device) AddDropListener(channel chan interface{}) {
	d.listeners[4] = append(d.listeners[4], channel)
}

func data_device_drop(d *Data_device, msg []byte) {
	var data Data_deviceDrop

	for _,channel := range d.listeners[4] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "drop", )
}

type Data_deviceSelection struct {
	Id *Data_offer
}

func (d *Data_device) AddSelectionListener(channel chan interface{}) {
	d.listeners[5] = append(d.listeners[5], channel)
}

func data_device_selection(d *Data_device, msg []byte) {
	var data Data_deviceSelection
	buf := bytes.NewBuffer(msg)

	idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(Data_offer)
	idobj := getObject(idid)
	if idobj == nil {
		return
	}
	id = idobj.(*Data_offer)
	data.Id = id

	for _,channel := range d.listeners[5] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "selection", id)
}

func NewData_device() (d *Data_device) {
	d = new(Data_device)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, data_device_data_offer)
	d.events = append(d.events, data_device_enter)
	d.events = append(d.events, data_device_leave)
	d.events = append(d.events, data_device_motion)
	d.events = append(d.events, data_device_drop)
	d.events = append(d.events, data_device_selection)
	return
}

func (d *Data_device) SetId(id int32) {
	d.id = id
}

func (d *Data_device) Id() int32 {
	return d.id
}