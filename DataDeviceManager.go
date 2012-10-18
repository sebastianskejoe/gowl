package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type DataDeviceManager struct {
	id int32
    
	events []func(d *DataDeviceManager, msg message) error
    name string
}

func NewDataDeviceManager() (d *DataDeviceManager) {
	d = new(DataDeviceManager)
    d.name = "DataDeviceManager"
    

    
	return
}

func (d *DataDeviceManager) HandleEvent(msg message) {
	if d.events[msg.opcode] != nil {
		d.events[msg.opcode](d, msg)
	}
}

func (d *DataDeviceManager) SetId(id int32) {
	d.id = id
}

func (d *DataDeviceManager) Id() int32 {
	return d.id
}

func (d *DataDeviceManager) Name() string {
    return d.name
}

////
//// REQUESTS
////

func (d *DataDeviceManager) CreateDataSource(id *DataSource) {
    sendrequest(d, "wl_data_device_manager_create_data_source", id)
}

func (d *DataDeviceManager) GetDataDevice(id *DataDevice, seat *Seat) {
    sendrequest(d, "wl_data_device_manager_get_data_device", id, seat)
}

////
//// EVENTS
////
