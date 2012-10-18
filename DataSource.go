package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataSource struct {
	id int32
    targetListeners []chan DataSourceTarget
	sendListeners []chan DataSourceSend
	cancelledListeners []chan DataSourceCancelled
	events []func(d *DataSource, msg message) error
    name string
}

func NewDataSource() (d *DataSource) {
	d = new(DataSource)
    d.name = "DataSource"
    d.targetListeners = make([]chan DataSourceTarget, 0)
	d.sendListeners = make([]chan DataSourceSend, 0)
	d.cancelledListeners = make([]chan DataSourceCancelled, 0)

    d.events = append(d.events, dataSourceTarget)
	d.events = append(d.events, dataSourceSend)
	d.events = append(d.events, dataSourceCancelled)
	return
}

func (d *DataSource) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

func (d *DataSource) SetId(id int32) {
	d.id = id
}

func (d *DataSource) Id() int32 {
	return d.id
}

func (d *DataSource) Name() string {
    return d.name
}

////
//// REQUESTS
////

func (d *DataSource) Offer(typ string) {
    sendrequest(d, "wl_data_source_offer", typ)
}

func (d *DataSource) Destroy() {
    sendrequest(d, "wl_data_source_destroy", )
}

////
//// EVENTS
////

type DataSourceTarget struct {
    MimeType string
}

func (d *DataSource) AddTargetListener(channel chan DataSourceTarget) {
    d.targetListeners = append(d.targetListeners, channel)
}

func dataSourceTarget(d *DataSource, msg message) (err error) {
    var data DataSourceTarget

    // Read mime_type
    mime_type,err := readString(msg.buf)
    if err != nil {
        return
    }
    data.MimeType = mime_type

    // Dispatch events
    for _,channel := range d.targetListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DataSourceSend struct {
    MimeType string
	Fd uintptr
}

func (d *DataSource) AddSendListener(channel chan DataSourceSend) {
    d.sendListeners = append(d.sendListeners, channel)
}

func dataSourceSend(d *DataSource, msg message) (err error) {
    var data DataSourceSend

    // Read mime_type
    mime_type,err := readString(msg.buf)
    if err != nil {
        return
    }
    data.MimeType = mime_type

    // Read fd
    fd,err := msg.fd, nil
    if err != nil {
        return
    }
    data.Fd = fd

    // Dispatch events
    for _,channel := range d.sendListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}

type DataSourceCancelled struct {
    
}

func (d *DataSource) AddCancelledListener(channel chan DataSourceCancelled) {
    d.cancelledListeners = append(d.cancelledListeners, channel)
}

func dataSourceCancelled(d *DataSource, msg message) (err error) {
    var data DataSourceCancelled


    // Dispatch events
    for _,channel := range d.cancelledListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
