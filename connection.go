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
    alive bool
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
    conn.alive = true

	return nil
}

func getmsg() (msgs []message, err error) {
    // Get all messages in buffer
    b := make([]byte, 1024)
    oob := make([]byte, 1024)
    n, oobn, flags,_, err := conn.unixconn.ReadMsgUnix(b, oob)
    fd := uintptr(0)

    if err != nil {
        return
    }

    bbuf := bytes.NewBuffer(b)
    oobbuf := bytes.NewBuffer(oob)

    if oobn != 0 {
        fmt.Println("Out-of-band recv: Len",oobn,oobbuf.Bytes()[:oobn], flags)
        readInt32(oobbuf)
        readInt32(oobbuf)
        readInt32(oobbuf)
        fd,err = readUintptr(oobbuf)
        fmt.Println("Fd is",fd)
        dup,_ := syscall.Dup(int(fd))
        fd = uintptr(dup)
        fmt.Println("Fd is",fd)
        if err != nil {
            fmt.Println("Had an error:", err)
            return
        }
    }


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

        if fd != 0 {
            msg.fd = fd
        }

        msgs = append(msgs, msg)

        fmt.Printf("%s@%d.%d(%d)\n", msg.obj.Name(), msg.obj.Id(), msg.opcode, msg.buf.Bytes())
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

func sendrequest(obj Object, req string, args ...interface{}) {
    if conn.alive != true {
        return
    }
    data := new(bytes.Buffer)
    dbg := new(bytes.Buffer)
    size := 8
    fd := uintptr(0)
    signature := signatures[req]
    opcode := signature.opcode
    sig := signature.signature
    for i,arg := range args {
        switch sig[i] {
        case 'i', 'f':
            size += 4
            binary.Write(data, binary.LittleEndian, arg.(int32))
            fmt.Fprintf(dbg, "%d ", arg)
        case 'u':
            size += 4
            binary.Write(data, binary.LittleEndian, arg.(uint32))
            fmt.Fprintf(dbg, "%d ", arg)
        case 'h':
	        fd,_,_ = syscall.Syscall(syscall.SYS_FCNTL, arg.(uintptr), syscall.F_DUPFD_CLOEXEC, 0)
        case 'o':
            size += 4
            binary.Write(data, binary.LittleEndian, arg.(Object).Id())
            fmt.Fprintf(dbg, "%d ", arg.(Object).Id())
        case 'n':
            size += 4
            nobj := arg.(Object)
            appendObject(nobj)
            binary.Write(data, binary.LittleEndian, nobj.Id())
            fmt.Fprintf(dbg, "new_id %d ", nobj.Id())
        case 's':
	        // First get padding
            str := arg.(string)
	        pad := 4-(len(str) % 4)
	        binary.Write(data, binary.LittleEndian, int32(len(str)+pad))
	        binary.Write(data, binary.LittleEndian, []byte(str))
	        for i := 0 ; i < pad ; i++ {
		        binary.Write(data, binary.LittleEndian, []byte{0})
	        }
            size += len(str)+pad+4
            fmt.Fprintf(dbg, "%s ", str)
        }
    }

    fmt.Printf(" -> %s@%d ( %s)\n",req,obj.Id(), dbg)

    // Send message
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, obj.Id())
	binary.Write(buf, binary.LittleEndian, opcode)
	binary.Write(buf, binary.LittleEndian, int16(size))
	binary.Write(buf, binary.LittleEndian, data.Bytes())

	var cmsgbytes []byte
	if fd != 0 {
		cmsgbytes = syscall.UnixRights(int(fd))
	}
	_,_,err := conn.unixconn.WriteMsgUnix(buf.Bytes(), cmsgbytes, nil)
	if err != nil {
		fmt.Println("sendrequest",err)
        os.Exit(1)
		return
	}
	syscall.Close(int(fd))
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

func readString(c io.Reader) (string, error) {
	// First get string length
	strlen, err := readUint32(c)
	if err != nil {
		return "", err
	}
	// Now get string
	str := make([]byte, strlen-1)
	_,err = c.Read(str)
	if err != nil {
		return "", err
	}

	pad := 4-(strlen % 4)
	if pad == 4 {
		pad = 0
	}
	for i := uint32(0) ; i <= pad ; i++ {
		binary.Read(c, binary.LittleEndian, []byte{0})
	}
	return string(str), nil
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
