package main

import (
	"gowl"
	"time"
	"fmt"
	"strings"
//	"os"
	"syscall"
)

/*
#define _GNU_SOURCE
#include <stdlib.h>
*/
import "C"

type Display struct {
	display *gowl.Display
	compositor *gowl.Compositor
	shm *gowl.Shm
	shell *gowl.Shell
	pool *gowl.Shm_pool
	buffer *gowl.Buffer
	surface *gowl.Surface
	shell_surface *gowl.Shell_surface
	data []byte
}

var (
	col uint8
	add int8
)

func (d *Display) globalListener(c chan interface{}) {
	for e := range c {
		glob := e.(gowl.DisplayGlobal)
		switch strings.TrimSpace(glob.Iface) {
		case "wl_shell":
			d.display.Bind(glob.Name, glob.Iface, glob.Version, d.shell)
		case "wl_shm":
			d.display.Bind(glob.Name, glob.Iface, glob.Version, d.shm)
		case "wl_compositor":
			d.display.Bind(glob.Name, glob.Iface, glob.Version, d.compositor)
		}
	}
}

func main() {
	display := new(Display)
	display.display = gowl.NewDisplay()
	display.compositor = gowl.NewCompositor()
	display.shm = gowl.NewShm()
	display.shell = gowl.NewShell()
	display.pool = gowl.NewShm_pool()
	display.buffer = gowl.NewBuffer()
	display.surface = gowl.NewSurface()
	display.shell_surface = gowl.NewShell_surface()

	globchan := make(chan interface{})
	go display.globalListener(globchan)
	display.display.AddGlobalListener(globchan)

	display.display.Iterate()

	// Sync
	WaitForSync(display.display)

	// create pool
	fd := create_tmp()
	mmap,err := syscall.Mmap(int(fd), 0, 250000, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println(err)
	}
	display.data = mmap
	col = 0
	add = 1
	//syscall.CloseOnExec(fd)
	display.shm.Create_pool(display.pool, fd, 2500000)
	display.pool.Create_buffer(display.buffer, 0, 250, 250, 1000, 1)
	display.pool.Destroy()

	// Create surfaces
	display.compositor.Create_surface(display.surface)
	display.shell.Get_shell_surface(display.shell_surface, display.surface)
	go Pong(display.shell_surface)
	display.shell_surface.Set_toplevel()

	redraw(display)

	<-time.After(2e9)
	display.buffer.Destroy()
	display.surface.Destroy()
}

func redraw(display *Display) {
	col = uint8(int8(col)+add)
	if col == 255 {
		add = -1
	} else if col == 0 {
		add = 1
	}

	for i,_ := range display.data {
		display.data[i] = byte(col)
	}
	display.surface.Attach(display.buffer, 0, 0)
	display.surface.Damage(0,0,250,250)
	cb := gowl.NewCallback()
	done := make(chan interface{})
	cb.AddDoneListener(done)
	display.surface.Frame(cb)
	func () {
		for {
			select {
			case <-done:
				redraw(display)
			default:
				display.display.Iterate()
			}
		}
	} ()
}

func WaitForSync(display *gowl.Display) {
	cb := gowl.NewCallback()
	done := make(chan interface{})
	cb.AddDoneListener(done)
	display.Sync(cb)
	func () {
		for {
			select {
			case <-done:
				return
			default:
				display.Iterate()
			}
		}
	} ()
}

func create_tmp() (uintptr) {
/*	_,_ = os.Create("bar")
//	file,err := os.Create("foo")//, os.O_EXCL|os.O_CREATE, 0777)
	file,err := os.OpenFile("/home/sebastian/.weston-tmp/foo", os.O_EXCL|os.O_CREATE|os.O_RDWR|syscall.O_CLOEXEC, 0600)
	os.Remove("/home/sebastian/.weston-tmp/foo")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	err = file.Truncate(250000)
	if err != nil {
		fmt.Println(err, file.Fd())
		return 0
	}
	fd,_,_ := syscall.Syscall(syscall.SYS_FCNTL, file.Fd(), syscall.F_DUPFD_CLOEXEC, 0)
	file.Close()

	return fd*/
	name := C.CString("/home/sebastian/.weston-tmp/gowl-XXXXXX")
	fd := uintptr(C.mkostemp(name, syscall.O_CLOEXEC))
	syscall.Ftruncate(int(fd), 250000)
	syscall.Unlink(C.GoString(name))
	return fd
}

func Pong(ss *gowl.Shell_surface) {
	c := make(chan interface{})
	ss.AddPingListener(c)
	for p := range c {
		ping := p.(gowl.Shell_surfacePing)
		ss.Pong(ping.Serial)
	}
}
