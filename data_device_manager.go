
package gowl

import (
	"bytes"
)

type Data_device_manager struct {
	*WlObject
	events []func (d *Data_device_manager, msg []byte)
}

//// Requests
func (d *Data_device_manager) Create_data_source (id *Data_source ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())

	sendmsg(d, 0, buf.Bytes())
}

func (d *Data_device_manager) Get_data_device (id *Data_device, seat *Seat ) {
	buf := new(bytes.Buffer)
	appendObject(id)
	writeInteger(buf, id.Id())
	writeInteger(buf, seat.Id())

	sendmsg(d, 1, buf.Bytes())
}

//// Events
func (d *Data_device_manager) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

func NewData_device_manager() (d *Data_device_manager) {
	d = new(Data_device_manager)

	return
}