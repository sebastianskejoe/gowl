package main

import (
	"gowl"
	"fmt"
	"strings"
	"syscall"
	"io/ioutil"
	"os"
)

type Display struct {
	display *gowl.Display
	compositor *gowl.Compositor
	shm *gowl.Shm
	shell *gowl.Shell
	pool *gowl.ShmPool
	buffer *gowl.Buffer
	surface *gowl.Surface
	shell_surface *gowl.ShellSurface
	data []byte
}

var (
	col uint8
	add int8
)

func main() {
	display := new(Display)
	display.display = gowl.NewDisplay()
	display.compositor = gowl.NewCompositor()
	display.shm = gowl.NewShm()
	display.shell = gowl.NewShell()
	display.pool = gowl.NewShmPool()
	display.buffer = gowl.NewBuffer()
	display.surface = gowl.NewSurface()
	display.shell_surface = gowl.NewShellSurface()

	globchan := make(chan interface{})
	go display.globalListener(globchan)
	display.display.AddGlobalListener(globchan)

	display.display.Iterate()

	// Sync
	waitForSync(display.display)

	// create pool
	fd := create_tmp(250*250*4)
	mmap,err := syscall.Mmap(int(fd), 0, 250000, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println(err)
	}
	display.data = mmap
	col = 0
	add = 1
	display.shm.CreatePool(display.pool, fd, 2500000)
	display.pool.CreateBuffer(display.buffer, 0, 250, 250, 1000, 1)
	display.pool.Destroy()

	// Create surfaces
	display.compositor.CreateSurface(display.surface)
	display.shell.GetShellSurface(display.shell_surface, display.surface)
	go Pong(display.shell_surface)
	display.shell_surface.SetToplevel()
	display.shell_surface.SetTitle("Gowl test window")

	redraw(display)

	display.buffer.Destroy()
	display.surface.Destroy()
}


//// Event listeners
func Pong(ss *gowl.ShellSurface) {
	c := make(chan interface{})
	ss.AddPingListener(c)
	for p := range c {
		ping := p.(gowl.ShellSurfacePing)
		ss.Pong(ping.Serial)
	}
}

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

//// Helper
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

func waitForSync(display *gowl.Display) {
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

func create_tmp(size int64) (uintptr) {
	tmp,err := ioutil.TempFile(os.Getenv("XDG_RUNTIME_DIR"), "gowl")
	if err != nil {
		fmt.Println(err)
	}
	tmp.Truncate(size)
	os.Remove(tmp.Name())
	return tmp.Fd()
}
