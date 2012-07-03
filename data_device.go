
package gowl

import (
	"bytes"
)

type Data_device struct {
	*WlObject
	events []func (d *Data_device, msg []byte)
}

//// Requests
func (d *Data_device) Start_drag (source *Data_source, origin *Surface, icon *Surface, serial uint32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, source.Id())
	writeInteger(buf, origin.Id())
	writeInteger(buf, icon.Id())
	writeInteger(buf, serial)

	sendmsg(d, 0, buf.Bytes())
}

func (d *Data_device) Set_selection (source *Data_source, serial uint32 ) {
	buf := new(bytes.Buffer)
	writeInteger(buf, source.Id())
	writeInteger(buf, serial)

	sendmsg(d, 1, buf.Bytes())
}

//// Events
func (d *Data_device) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

type Data_deviceData_offer struct {
	id *Data_offer
}

func (d *Data_device) AddData_offerListener(channel chan interface{}) {
	d.addListener(0, channel)
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
	data.id = id

	for _,channel := range d.listeners[0] {
		channel <- data
	}
}

type Data_deviceEnter struct {
	serial uint32
	surface *Surface
	x int32
	y int32
	id *Data_offer
}

func (d *Data_device) AddEnterListener(channel chan interface{}) {
	d.addListener(1, channel)
}

func data_device_enter(d *Data_device, msg []byte) {
	var data Data_deviceEnter
	buf := bytes.NewBuffer(msg)

	serial,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.serial = serial

	surfaceid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	surface := new(Surface)
	surface = getObject(surfaceid).(*Surface)
	data.surface = surface

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

	idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(Data_offer)
	id = getObject(idid).(*Data_offer)
	data.id = id

	for _,channel := range d.listeners[1] {
		channel <- data
	}
}

type Data_deviceLeave struct {
}

func (d *Data_device) AddLeaveListener(channel chan interface{}) {
	d.addListener(2, channel)
}

func data_device_leave(d *Data_device, msg []byte) {
	var data Data_deviceLeave

	for _,channel := range d.listeners[2] {
		channel <- data
	}
}

type Data_deviceMotion struct {
	time uint32
	x int32
	y int32
}

func (d *Data_device) AddMotionListener(channel chan interface{}) {
	d.addListener(3, channel)
}

func data_device_motion(d *Data_device, msg []byte) {
	var data Data_deviceMotion
	buf := bytes.NewBuffer(msg)

	time,err := readUint32(buf)
	if err != nil {
		// XXX Error handling
	}
	data.time = time

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

	for _,channel := range d.listeners[3] {
		channel <- data
	}
}

type Data_deviceDrop struct {
}

func (d *Data_device) AddDropListener(channel chan interface{}) {
	d.addListener(4, channel)
}

func data_device_drop(d *Data_device, msg []byte) {
	var data Data_deviceDrop

	for _,channel := range d.listeners[4] {
		channel <- data
	}
}

type Data_deviceSelection struct {
	id *Data_offer
}

func (d *Data_device) AddSelectionListener(channel chan interface{}) {
	d.addListener(5, channel)
}

func data_device_selection(d *Data_device, msg []byte) {
	var data Data_deviceSelection
	buf := bytes.NewBuffer(msg)

	idid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	id := new(Data_offer)
	id = getObject(idid).(*Data_offer)
	data.id = id

	for _,channel := range d.listeners[5] {
		channel <- data
	}
}

func NewData_device() (d *Data_device) {
	d = new(Data_device)

	d.events = append(d.events, data_device_data_offer)
	d.events = append(d.events, data_device_enter)
	d.events = append(d.events, data_device_leave)
	d.events = append(d.events, data_device_motion)
	d.events = append(d.events, data_device_drop)
	d.events = append(d.events, data_device_selection)
	return
}