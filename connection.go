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
	unixconn *net.UnixConn
	addr *net.UnixAddr
	reader *bufio.Reader
	writer *bufio.Writer
}

type message struct {
	obj Object
	opcode int16
	buf *bytes.Buffer
	fd uintptr
}

var conn Connection

func connect_to_socket() error {
	// Connect to socket
	addr := fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), "wayland-0")
	c,err := net.DialUnix("unix", conn.addr, &net.UnixAddr{addr, "unix"})
	if err != nil {
		return err
	}

	conn.reader = bufio.NewReader(c)
	conn.writer = bufio.NewWriter(c)
	conn.unixconn = c

	return nil
}

func getmsg() (msgs []message, err error) {
    // Get all messages in buffer
    b := make([]byte, 1024)
    oob := make([]byte, 1024)
    n, oobn, _,_, err := conn.unixconn.ReadMsgUnix(b, oob)

    if err != nil {
        return
    }

    bbuf := bytes.NewBuffer(b)
    oobbuf := bytes.NewBuffer(oob)
    read := 0
    for read < n {
        var msg message

        id,err := readInt32(bbuf)
        if err != nil {
            break
        }
        msg.obj = getObject(id)

        msg.opcode, err = readInt16(bbuf)
        if err != nil {
            break
        }

        size, err := readInt16(bbuf)
        if err != nil {
            break
        }

        msg.buf = bytes.NewBuffer(bbuf.Next(int(size-8)))
        if err != nil {
            break
        }
        read += int(size)

        if oobn != 0 {
            msg.fd,err = readUintptr(oobbuf)
            if err != nil {
                fmt.Println("Had an error:", err)
                break
            }
        }
        msgs = append(msgs, msg)
    }

	return
}

func sendmsg(msg *message) {
	size := len(msg.buf.Bytes())
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, msg.obj.Id())
	binary.Write(buf, binary.LittleEndian, msg.opcode)
	binary.Write(buf, binary.LittleEndian, int16(size+8))
	binary.Write(buf, binary.LittleEndian, msg.buf.Bytes())

	var cmsgbytes []byte
	if msg.fd != 0 {
		cmsgbytes = syscall.UnixRights(int(msg.fd))
	}
	_,_,err := conn.unixconn.WriteMsgUnix(buf.Bytes(), cmsgbytes, nil)
	if err != nil {
		fmt.Println("sendmsg",err)
		return
	}
	syscall.Close(int(msg.fd))
}

func readUintptr(buf io.Reader) (uintptr, error) {
	var val uint32
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		return uintptr(val), err
	}
	return uintptr(val), nil
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

func readFixed(c io.Reader) (int32, error) {
	var val int32
	err := binary.Read(c, binary.LittleEndian, &val)
	if err != nil {
		return val, err
	}
	return val/256, nil
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
	str := make([]byte, strlen-1)
	_,err = c.Read(str)
	if err != nil {
		return 0, "", err
	}

	pad := 4-(strlen % 4)
	if pad == 4 {
		pad = 0
	}
	for i := uint32(0) ; i <= pad ; i++ {
		binary.Read(c, binary.LittleEndian, []byte{0})
	}
	return strlen, string(str), nil
}

func readArray(c io.Reader) ([]interface{}, error) {
	return nil, nil
}

func writeInteger(msg *message, val interface{}) {
	err := binary.Write(msg.buf, binary.LittleEndian, val)
	if err != nil {
		fmt.Println(err)
	}
}

func writeString(msg *message, val []byte) {
	// First get padding
	pad := 4-(len(val) % 4)

	writeInteger(msg, int32(len(val)+pad))
	binary.Write(msg.buf, binary.LittleEndian, val)
	for i := 0 ; i < pad ; i++ {
		binary.Write(msg.buf, binary.LittleEndian, []byte{0})
	}
}

func writeFd(msg *message, val uintptr) {
	newfd,_,_ := syscall.Syscall(syscall.SYS_FCNTL, val, syscall.F_DUPFD_CLOEXEC, 0)
	msg.fd = newfd
}

func newMessage(obj Object, opcode int16) *message {
	return &message{
		obj: obj,
		opcode: opcode,
		buf: new(bytes.Buffer),
		fd: 0,
	}
}
