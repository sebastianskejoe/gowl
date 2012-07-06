package gowl
type Data_device_manager struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (d *Data_device_manager, msg []byte)
}

//// Requests
func (d *Data_device_manager) Create_data_source (id *Data_source) {
	msg := newMessage(d, 0)
	appendObject(id)
	writeInteger(msg,id.Id())

	sendmsg(msg)
	printRequest("data_device_manager", "create_data_source", id)
}

func (d *Data_device_manager) Get_data_device (id *Data_device, seat *Seat) {
	msg := newMessage(d, 1)
	appendObject(id)
	writeInteger(msg,id.Id())
	writeInteger(msg,seat.Id())

	sendmsg(msg)
	printRequest("data_device_manager", "get_data_device", id, seat)
}

//// Events
func (d *Data_device_manager) HandleEvent(opcode int16, msg []byte) {
	if d.events[opcode] != nil {
		d.events[opcode](d, msg)
	}
}

func NewData_device_manager() (d *Data_device_manager) {
	d = new(Data_device_manager)
	d.listeners = make(map[int16][]chan interface{}, 0)

	return
}

func (d *Data_device_manager) SetId(id int32) {
	d.id = id
}

func (d *Data_device_manager) Id() int32 {
	return d.id
}