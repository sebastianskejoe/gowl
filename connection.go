package gowl

import (
	"encoding/binary"
	"bytes"
	"net"
	"fmt"
	"io"
)

func sendmsg(conn net.Conn, id int32, opcode int16, msg []byte) {
	size := len(msg)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, id)
	binary.Write(buf, binary.LittleEndian, opcode)
	binary.Write(buf, binary.LittleEndian, int16(size+8))
	binary.Write(buf, binary.LittleEndian, msg)
	fmt.Printf("Sending message %d\n", buf.Bytes())
	fmt.Fprintf(conn, "%s", buf.Bytes())
}

func readUint32(buf io.Reader) (uint32) {
	var val uint32
	binary.Read(buf, binary.LittleEndian, &val)
	return val
}

func readInt32(buf io.Reader) (int32) {
	var val int32
	fmt.Println("Waiting for read")
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		fmt.Println(err)
	}
	return val
}

func readUint16(buf io.Reader) (uint16) {
	var val uint16
	binary.Read(buf, binary.LittleEndian, &val)
	return val
}

func readInt16(buf io.Reader) (int16) {
	var val int16
	binary.Read(buf, binary.LittleEndian, &val)
	return val
}

func readString(buf io.Reader) (uint32, []byte) {
	// First get string length
	strlen := readUint32(buf)
	// Now get string
	str := make([]byte, strlen)
	buf.Read(str)

	return strlen, str
}

func writeInteger(buf io.Writer, val interface{}) {
	binary.Write(buf, binary.LittleEndian, val)
}

func writeString(buf io.Writer, val []byte) {
	writeInteger(buf, int32(len(val)))
	binary.Write(buf, binary.LittleEndian, val)
	pad := len(val) % 4
	for i := 0 ; i < pad ; i++ {
		binary.Write(buf, binary.LittleEndian, '0')
	}
}
