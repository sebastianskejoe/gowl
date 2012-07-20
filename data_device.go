package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataDevice struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *DataDevice, msg []byte)
}

//// Requests
func (d *DataDevice) StartDrag (source *DataSource, origin *Surface, icon *Surface, serial uint32) {
	msg := newMessage(d, 0)
	writeInteger(msg,source.Id())
	writeInteger(msg,origin.Id())
	writeInteger(msg,icon.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("data_device", "start_drag", source, origin, icon, serial)
}

func (d *DataDevice) SetSelection (source *DataSource, serial uint32) {
	msg := newMessage(d, 1)
	writeInteger(msg,source.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("data_device", "set_selection", source, serial)
}

//// Events
func (d *DataDevice) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type DataDeviceDataOffer struct {
	Id *DataOffer
}

func (d *DataDevice) AddDataOfferListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
}

func data_device_data_offer(d *DataDevice, msg []byte) {
	var data DataDeviceDataOffer
	buf := bytes.NewBuffer(msg)

	idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(DataOffer)
	setObject(idid, id)
	data.Id = id

	for _,channel := range d.listeners[0] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "data_offer", id)
}

type DataDeviceEnter struct {
	Serial uint32
	Surface *Surface
	X int32
	Y int32
	Id *DataOffer
}

func (d *DataDevice) AddEnterListener(channel chan interface{}) {
	d.listeners[1] = append(d.listeners[1], channel)
}

func data_device_enter(d *DataDevice, msg []byte) {
	var data DataDeviceEnter
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
	id := new(DataOffer)
	idobj := getObject(idid)
	if idobj == nil {
		return
	}
	id = idobj.(*DataOffer)
	data.Id = id

	for _,channel := range d.listeners[1] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "enter", serial, surface, x, y, id)
}

type DataDeviceLeave struct {
}

func (d *DataDevice) AddLeaveListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
}

func data_device_leave(d *DataDevice, msg []byte) {
	var data DataDeviceLeave

	for _,channel := range d.listeners[2] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "leave", )
}

type DataDeviceMotion struct {
	Time uint32
	X int32
	Y int32
}

func (d *DataDevice) AddMotionListener(channel chan interface{}) {
	d.listeners[3] = append(d.listeners[3], channel)
}

func data_device_motion(d *DataDevice, msg []byte) {
	var data DataDeviceMotion
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

type DataDeviceDrop struct {
}

func (d *DataDevice) AddDropListener(channel chan interface{}) {
	d.listeners[4] = append(d.listeners[4], channel)
}

func data_device_drop(d *DataDevice, msg []byte) {
	var data DataDeviceDrop

	for _,channel := range d.listeners[4] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "drop", )
}

type DataDeviceSelection struct {
	Id *DataOffer
}

func (d *DataDevice) AddSelectionListener(channel chan interface{}) {
	d.listeners[5] = append(d.listeners[5], channel)
}

func data_device_selection(d *DataDevice, msg []byte) {
	var data DataDeviceSelection
	buf := bytes.NewBuffer(msg)

	idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(DataOffer)
	idobj := getObject(idid)
	if idobj == nil {
		return
	}
	id = idobj.(*DataOffer)
	data.Id = id

	for _,channel := range d.listeners[5] {
		go func () { channel <- data }()
	}
	printEvent("data_device", "selection", id)
}

func NewDataDevice() (d *DataDevice) {
	d = new(DataDevice)
	d.listeners = make(map[int16][]chan interface{}, 0)

	d.events = append(d.events, data_device_data_offer)
	d.events = append(d.events, data_device_enter)
	d.events = append(d.events, data_device_leave)
	d.events = append(d.events, data_device_motion)
	d.events = append(d.events, data_device_drop)
	d.events = append(d.events, data_device_selection)
	return
}

func (d *DataDevice) SetId(id int32) {
	d.id = id
}

func (d *DataDevice) Id() int32 {
	return d.id
}