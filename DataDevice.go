package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataDevice struct {
	id int32
    dataOfferListeners []chan DataDeviceDataOffer
	enterListeners []chan DataDeviceEnter
	leaveListeners []chan DataDeviceLeave
	motionListeners []chan DataDeviceMotion
	dropListeners []chan DataDeviceDrop
	selectionListeners []chan DataDeviceSelection
	events []func(d *DataDevice, msg message) error
    name string
}

func NewDataDevice() (d *DataDevice) {
	d = new(DataDevice)
    d.name = "DataDevice"
    d.dataOfferListeners = make([]chan DataDeviceDataOffer, 0)
	d.enterListeners = make([]chan DataDeviceEnter, 0)
	d.leaveListeners = make([]chan DataDeviceLeave, 0)
	d.motionListeners = make([]chan DataDeviceMotion, 0)
	d.dropListeners = make([]chan DataDeviceDrop, 0)
	d.selectionListeners = make([]chan DataDeviceSelection, 0)

    d.events = append(d.events, dataDeviceDataOffer)
	d.events = append(d.events, dataDeviceEnter)
	d.events = append(d.events, dataDeviceLeave)
	d.events = append(d.events, dataDeviceMotion)
	d.events = append(d.events, dataDeviceDrop)
	d.events = append(d.events, dataDeviceSelection)
	return
}

func (d *DataDevice) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

func (d *DataDevice) SetId(id int32) {
	d.id = id
}

func (d *DataDevice) Id() int32 {
	return d.id
}

func (d *DataDevice) Name() string {
    return d.name
}

////
//// REQUESTS
////

func (d *DataDevice) StartDrag(source *DataSource, origin *Surface, icon *Surface, serial uint32) {
    sendrequest(d, "wl_data_device_start_drag", source, origin, icon, serial)
}

func (d *DataDevice) SetSelection(source *DataSource, serial uint32) {
    sendrequest(d, "wl_data_device_set_selection", source, serial)
}

////
//// EVENTS
////

type DataDeviceDataOffer struct {
    Id *DataOffer
}

func (d *DataDevice) AddDataOfferListener(channel chan DataDeviceDataOffer) {
    d.dataOfferListeners = append(d.dataOfferListeners, channel)
}

func dataDeviceDataOffer(d *DataDevice, msg message) (err error) {
    var data DataDeviceDataOffer

    // Read id
    id,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    idObj := new(DataOffer)
    setObject(id, idObj)
    data.Id = idObj

    // Dispatch events
    for _,channel := range d.dataOfferListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DataDeviceEnter struct {
    Serial uint32
	Surface *Surface
	X int32
	Y int32
	Id *DataOffer
}

func (d *DataDevice) AddEnterListener(channel chan DataDeviceEnter) {
    d.enterListeners = append(d.enterListeners, channel)
}

func dataDeviceEnter(d *DataDevice, msg message) (err error) {
    var data DataDeviceEnter

    // Read serial
    serial,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Serial = serial

    // Read surface
    surface,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    surfaceObj := getObject(surface)
    data.Surface = surfaceObj.(*Surface)

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

    // Read id
    id,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    idObj := getObject(id)
    data.Id = idObj.(*DataOffer)

    // Dispatch events
    for _,channel := range d.enterListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DataDeviceLeave struct {
    
}

func (d *DataDevice) AddLeaveListener(channel chan DataDeviceLeave) {
    d.leaveListeners = append(d.leaveListeners, channel)
}

func dataDeviceLeave(d *DataDevice, msg message) (err error) {
    var data DataDeviceLeave


    // Dispatch events
    for _,channel := range d.leaveListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DataDeviceMotion struct {
    Time uint32
	X int32
	Y int32
}

func (d *DataDevice) AddMotionListener(channel chan DataDeviceMotion) {
    d.motionListeners = append(d.motionListeners, channel)
}

func dataDeviceMotion(d *DataDevice, msg message) (err error) {
    var data DataDeviceMotion

    // Read time
    time,err := readUint32(msg.buf)
    if err != nil {
        return
    }
    data.Time = time

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
    for _,channel := range d.motionListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DataDeviceDrop struct {
    
}

func (d *DataDevice) AddDropListener(channel chan DataDeviceDrop) {
    d.dropListeners = append(d.dropListeners, channel)
}

func dataDeviceDrop(d *DataDevice, msg message) (err error) {
    var data DataDeviceDrop


    // Dispatch events
    for _,channel := range d.dropListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DataDeviceSelection struct {
    Id *DataOffer
}

func (d *DataDevice) AddSelectionListener(channel chan DataDeviceSelection) {
    d.selectionListeners = append(d.selectionListeners, channel)
}

func dataDeviceSelection(d *DataDevice, msg message) (err error) {
    var data DataDeviceSelection

    // Read id
    id,err := readInt32(msg.buf)
    if err != nil {
        return
    }
    idObj := getObject(id)
    data.Id = idObj.(*DataOffer)

    // Dispatch events
    for _,channel := range d.selectionListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
