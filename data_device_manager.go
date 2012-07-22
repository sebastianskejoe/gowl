package gowl

type DataDeviceManager struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *DataDeviceManager, msg []byte)
}

//// Requests
func (d *DataDeviceManager) CreateDataSource (id *DataSource) {
	msg := newMessage(d, 0)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("data_device_manager", d, "create_data_source", "new id", id.Id())
}

func (d *DataDeviceManager) GetDataDevice (id *DataDevice, seat *Seat) {
	msg := newMessage(d, 1)
	appendObject(id)
	writeInteger(msg,id.Id())
	writeInteger(msg,seat.Id())

	sendmsg(msg)
	printRequest("data_device_manager", d, "get_data_device", "new id", id.Id(), seat.Id())
}

//// Events
func (d *DataDeviceManager) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

func NewDataDeviceManager() (d *DataDeviceManager) {
	d = new(DataDeviceManager)
	d.listeners = make(map[int16][]chan interface{}, 0)

	return
}

func (d *DataDeviceManager) SetId(id int32) {
	d.id = id
}

func (d *DataDeviceManager) Id() int32 {
	return d.id
}