package gowl

import (
	"net"
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"os"
	"bufio"
)

type Display struct {
	events []func (d *Display, msg []byte)
	conn net.Conn
	id int32
}

func NewDisplay() *Display {
	// Connect to socket
	addr := fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), "wayland-0")
	conn, err := net.Dial("unix", addr)
	if err != nil {
		fmt.Println("Couldn't connect", err)
		return nil
	}

	// Create display struct
	d := new(Display)
	d.conn = conn
	d.events = append(d.events, display_error)
	d.events = append(d.events, display_global)

	appendObject(d)

	go func () {
	r := bufio.NewReader(conn)
	for {
		fmt.Println("Waiting for event")
		// Message header
		id := readInt32(r)
		fmt.Println("Waiting for read")
		opcode := readInt16(r)
		size := readInt16(r)

		// Message
		msg := make([]byte, size-8)
		r.Read(msg)

		if (objects[id] != nil) {
			objects[id].HandleEvent(opcode, msg)
		}
	}
	}()

	return d
}

func (d *Display) SetID(id int32) {
	d.id = id
}

func (d *Display) ID() int32 {
	return d.id
}

func (d *Display) HandleEvent(opcode int16, msg []byte) {
	d.events[opcode](d,msg)
}

func (d *Display) Bind(name uint32, iface []byte, version uint32, obj Object) (int32) {
	appendObject(obj)
	// Create message
	msg := new(bytes.Buffer)
	writeInteger(msg, name)
	writeString(msg, iface)
	writeInteger(msg, version)
	binary.Write(msg, binary.LittleEndian, obj.ID())

	sendmsg(d.conn, 1, 0, msg.Bytes())
	return obj.ID()
}

func display_error(d *Display, msg []byte) {
	buf := bytes.NewBuffer(msg)
	object := readInt32(buf)
	err := readUint32(buf)
	_, message := readString(buf)
	if (object != 0) {
		fmt.Printf("ERROR { %d, %d, %s }\n", object, err, message)
	}
}

func display_global(d *Display, msg []byte) {
	buf := bytes.NewBuffer(msg)
	name := readUint32(buf)
	_, iface := readString(buf)
	version := readUint32(buf)

	if name == 0 {
		return
	}

	fmt.Printf("GLOBAL { %d, %s, %d }\n", name, iface, version)

	if strings.Contains(string(iface), "wl_compositor") {
		fmt.Printf("ID: %d\n", d.Bind(name, iface, version, nil))
	} else if strings.Contains(string(iface), "wl_shm") {
		fmt.Printf("ID: %d\n", d.Bind(name, iface, version, nil))
	} else if strings.Contains(string(iface), "wl_shell") {
		fmt.Printf("ID: %d\n", d.Bind(name, iface, version, nil))
	}
}
