package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataDevice struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *DataDevice, msg message)
}

//// Requests
func (d *DataDevice) StartDrag (source *DataSource, origin *Surface, icon *Surface, serial uint32) {
	msg := newMessage(d, 0)
	writeInteger(msg,source.Id())
	writeInteger(msg,origin.Id())
	writeInteger(msg,icon.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("data_device", d, "start_drag", source.Id(), origin.Id(), icon.Id(), serial)
}

func (d *DataDevice) SetSelection (source *DataSource, serial uint32) {
	msg := newMessage(d, 1)
	writeInteger(msg,source.Id())
	writeInteger(msg,serial)

	sendmsg(msg)
	printRequest("data_device", d, "set_selection", source.Id(), serial)
}

//// Events
func (d *DataDevice) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

type DataDeviceDataOffer struct {
	Id *DataOffer
}

func (d *DataDevice) AddDataOfferListener(channel chan interface{}) {
	d.listeners[0] = append(d.listeners[0], channel)
	addListener(channel)
}

func data_device_data_offer(d *DataDevice, msg message) {
	var data DataDeviceDataOffer

	idid, err := readInt32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(DataOffer)
	setObject(idid, id)
	data.Id = id

	for _,channel := range d.listeners[0] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("data_device", d, "data_offer", "new id", id.Id())
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
	addListener(channel)
}

func data_device_enter(d *DataDevice, msg message) {
	var data DataDeviceEnter

	serial,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Serial = serial

	surfaceid, err := readInt32(msg.buf)
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

	x,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.X = x

	y,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Y = y

	idid, err := readInt32(msg.buf)
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
		go func() {
			channel <- data
		} ()
	}
	printEvent("data_device", d, "enter", serial, surface.Id(), x, y, id.Id())
}

type DataDeviceLeave struct {
}

func (d *DataDevice) AddLeaveListener(channel chan interface{}) {
	d.listeners[2] = append(d.listeners[2], channel)
	addListener(channel)
}

func data_device_leave(d *DataDevice, msg message) {
	var data DataDeviceLeave

	for _,channel := range d.listeners[2] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("data_device", d, "leave")
}

type DataDeviceMotion struct {
	Time uint32
	X int32
	Y int32
}

func (d *DataDevice) AddMotionListener(channel chan interface{}) {
	d.listeners[3] = append(d.listeners[3], channel)
	addListener(channel)
}

func data_device_motion(d *DataDevice, msg message) {
	var data DataDeviceMotion

	time,err := readUint32(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Time = time

	x,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.X = x

	y,err := readFixed(msg.buf)
	if err != nil {
		// XXX Error handling
	}
	data.Y = y

	for _,channel := range d.listeners[3] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("data_device", d, "motion", time, x, y)
}

type DataDeviceDrop struct {
}

func (d *DataDevice) AddDropListener(channel chan interface{}) {
	d.listeners[4] = append(d.listeners[4], channel)
	addListener(channel)
}

func data_device_drop(d *DataDevice, msg message) {
	var data DataDeviceDrop

	for _,channel := range d.listeners[4] {
		go func() {
			channel <- data
		} ()
	}
	printEvent("data_device", d, "drop")
}

type DataDeviceSelection struct {
	Id *DataOffer
}

func (d *DataDevice) AddSelectionListener(channel chan interface{}) {
	d.listeners[5] = append(d.listeners[5], channel)
	addListener(channel)
}

func data_device_selection(d *DataDevice, msg message) {
	var data DataDeviceSelection

	idid, err := readInt32(msg.buf)
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
		go func() {
			channel <- data
		} ()
	}
	printEvent("data_device", d, "selection", id.Id())
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