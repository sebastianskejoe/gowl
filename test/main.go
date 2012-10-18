package main

import (
	"fmt"
	"github.com/sebastianskejoe/gowl"
	"strings"
	"syscall"
    "time"
)

type Display struct {
	display       *gowl.Display
    registry      *gowl.Registry
	compositor    *gowl.Compositor
	shm           *gowl.Shm
	shell         *gowl.Shell
	pool          *gowl.ShmPool
	buffer        *gowl.Buffer
	surface       *gowl.Surface
	shell_surface *gowl.ShellSurface
	data          []byte
}

var (
	col uint8
	add int8
)

func main() {
	display := new(Display)
	display.display = gowl.NewDisplay()
    display.registry = gowl.NewRegistry()

	err := display.display.Connect()
	if err != nil {
		fmt.Println("Couldn't connect:", err)
		return
	}

	display.compositor = gowl.NewCompositor()
	display.shm = gowl.NewShm()
	display.shell = gowl.NewShell()
	display.pool = gowl.NewShmPool()
	display.buffer = gowl.NewBuffer()
	display.surface = gowl.NewSurface()
	display.shell_surface = gowl.NewShellSurface()

    display.display.GetRegistry(display.registry)
	globchan := make(chan gowl.RegistryGlobal)
	go display.globalListener(globchan)
	display.registry.AddGlobalListener(globchan)
	waitForSync(display.display)

	// create pool
	fd := gowl.CreateTmp(250 * 250 * 4)
	mmap, err := syscall.Mmap(int(fd), 0, 250000, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
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

	display.surface.Attach(display.buffer, 0, 0)
	redraw(display)

	display.buffer.Destroy()
	display.surface.Destroy()
}

//// Event listeners
func Pong(ss *gowl.ShellSurface) {
	c := make(chan gowl.ShellSurfacePing)
	ss.AddPingListener(c)
	for ping := range c {
		ss.Pong(ping.Serial)
	}
}

func (d *Display) globalListener(c chan gowl.RegistryGlobal) {
	for glob := range c {
		switch strings.TrimSpace(glob.Iface) {
		case "wl_shell":
            d.registry.Bind(glob.Name, glob.Iface, glob.Version, d.shell)
		case "wl_shm":
			d.registry.Bind(glob.Name, glob.Iface, glob.Version, d.shm)
		case "wl_compositor":
			d.registry.Bind(glob.Name, glob.Iface, glob.Version, d.compositor)
		}
	}
}

//// Helper
func redraw(display *Display) {
	col = uint8(int8(col) + add)
	if col == 255 {
		add = -1
	} else if col == 0 {
		add = 1
	}

	for i, _ := range display.data {
		display.data[i] = byte(col)
	}

	display.surface.Damage(0, 0, 250, 250)
	cb := gowl.NewCallback()
	done := make(chan gowl.CallbackDone)
	cb.AddDoneListener(done)
	display.surface.Frame(cb)
    display.surface.Commit()
	func() {
		for {
			select {
			case <-done:
				redraw(display)
			default:
                display.display.Iterate()
                c := time.Tick(time.Second/100)
                <-c
			}
		}
	}()
}

func waitForSync(display *gowl.Display) {
	cb := gowl.NewCallback()
	done := make(chan gowl.CallbackDone)
	cb.AddDoneListener(done)
	display.Sync(cb)
	func() {
		for {
			select {
			case <-done:
				return
			default:
				display.Iterate()
			}
		}
	}()
}
