package gowl

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"io"
	"bufio"
	"syscall"
	"os"
	"net"
)

type Connection struct {
	fd int
	sockaddr syscall.Sockaddr
	reader *bufio.Reader
	writer *bufio.Writer
}

var conn Connection

func connect_to_socket() {
	// Connect to socket
	addr := fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), "wayland-0")
	c,err := net.Dial("unix", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.reader = bufio.NewReader(c)
	conn.writer = bufio.NewWriter(c)
}

func getmsg() (id int32, opcode int16, size int16, msg []byte, remain int, err error) {
	// Message header
	id,err = readInt32(conn.reader)
	if err != nil {
		return
	}
	opcode,err = readInt16(conn.reader)
	if err != nil {
		return
//		fmt.Println(opcode, err)
//		break
	}
	size,err = readInt16(conn.reader)
	if err != nil {
		return
//		fmt.Println(size, err)
//		break
	}

	// Message
	msg = make([]byte, size-8)
	_,err = conn.reader.Read(msg)
	if err != nil {
		return
//		printError("getmsg", err)
//		break
	}

	remain = conn.reader.Buffered()

	return
}

func sendmsg(obj Object, opcode int16, msg []byte) {
	size := len(msg)
	buf := new(bytes.Buffer)
	writeInteger(buf, obj.Id())
	writeInteger(buf, opcode)
	writeInteger(buf, int16(size+8))
	binary.Write(buf, binary.LittleEndian, msg)

	err := binary.Write(conn.writer, binary.LittleEndian, buf.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	err = conn.writer.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readUintptr(buf io.Reader) (uintptr, error) {
	var val uintptr
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func readUint32(buf io.Reader) (uint32, error) {
	var val uint32
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func readInt32(c io.Reader) (int32, error) {
	var val int32
	err := binary.Read(c, binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func readUint16(c io.Reader) (uint16, error) {
	var val uint16
	err := binary.Read(c, binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func readInt16(c io.Reader) (int16, error) {
	var val int16
	err := binary.Read(c, binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func readString(c io.Reader) (uint32, string, error) {
	// First get string length
	strlen, err := readUint32(c)
	if err != nil {
		return 0, "", err
	}
	// Now get string
	str := make([]byte, strlen)
	_,err = c.Read(str)
	if err != nil {
		return 0, "", err
	}

	pad := 4-(strlen % 4)
	if pad == 4 {
		pad = 0
	}
	for i := uint32(0) ; i < pad ; i++ {
		binary.Read(c, binary.LittleEndian, []byte{0})
	}
	return strlen, string(str), nil
}

func readArray(c io.Reader) ([]interface{}, error) {
	return nil, nil
}

func writeInteger(c io.Writer, val interface{}) {
	binary.Write(c, binary.LittleEndian, val)
}

func writeString(c io.Writer, val []byte) {
	// First get padding
	pad := 4-(len(val) % 4)
	if pad == 4 {
		pad = 0
	}

	writeInteger(c, int32(len(val)+pad))
	binary.Write(c, binary.LittleEndian, val)
	for i := 0 ; i < pad ; i++ {
		binary.Write(c, binary.LittleEndian, []byte{0})
	}
}
