package main

import (
	"github.com/sebastianskejoe/gowl"

	"image"
	"image/draw"

	"fmt"
	"strings"
	"syscall"
)

type Window struct {
	display      *gowl.Display
	compositor   *gowl.Compositor
	surface      *gowl.Surface
	shell        *gowl.Shell
	shellsurface *gowl.ShellSurface
	shm          *gowl.Shm
	pool         *gowl.ShmPool
	buffer       *gowl.Buffer
	seat         *gowl.Seat
	pointer      *gowl.Pointer
	keyboard	 *gowl.Keyboard
	ddm          *gowl.DataDeviceManager
	dd           *gowl.DataDevice

	screen    *image.RGBA
	eventchan chan interface{}
}

func NewWindow(width, height int) (w *Window, err error) {
	w = new(Window)
	w.eventchan = make(chan interface{})

	// Create display and connect to wayland server
	w.display = gowl.NewDisplay()
	err = w.display.Connect()
	if err != nil {
		fmt.Println(err)
	}

	// Allocate other components
	w.compositor	= gowl.NewCompositor()
	w.surface		= gowl.NewSurface()
	w.shell			= gowl.NewShell()
	w.shellsurface	= gowl.NewShellSurface()

	w.shm		= gowl.NewShm()
	w.pool		= gowl.NewShmPool()
	w.buffer	= gowl.NewBuffer()

	w.seat		= gowl.NewSeat()
	w.pointer	= gowl.NewPointer()
	w.keyboard	= gowl.NewKeyboard()
	w.ddm		= gowl.NewDataDeviceManager()
	w.dd		= gowl.NewDataDevice()

	// Listen for global events from display
	globals := make(chan interface{})
	go func() {
		for event := range globals {
			global := event.(gowl.DisplayGlobal)
			switch strings.TrimSpace(global.Iface) {
			case "wl_compositor":
				w.display.Bind(global.Name, global.Iface, global.Version, w.compositor)
			case "wl_shm":
				w.display.Bind(global.Name, global.Iface, global.Version, w.shm)
			case "wl_shell":
				w.display.Bind(global.Name, global.Iface, global.Version, w.shell)
			case "wl_seat":
				w.display.Bind(global.Name, global.Iface, global.Version, w.seat)
				w.ddm.GetDataDevice(w.dd, w.seat)
				w.seat.GetPointer(w.pointer)
			case "wl_data_device_manager":
				w.display.Bind(global.Name, global.Iface, global.Version, w.ddm)
			}
		}
	}()
	w.display.AddGlobalListener(globals)

	// Iterate until we are sync'ed
	err = w.display.Iterate()
	if err != nil {
		w = nil
		return
	}
	waitForSync(w.display)

	// Create memory map
	w.screen = image.NewRGBA(image.Rect(0, 0, width, height))
	size := w.screen.Stride * w.screen.Rect.Dy()
	fd := gowl.CreateTmp(int64(size))
	mmap, err := syscall.Mmap(int(fd), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		w = nil
		return
	}
	w.screen.Pix = mmap

	// Create pool and buffer
	w.shm.CreatePool(w.pool, fd, int32(size))
	w.pool.CreateBuffer(w.buffer, 0, int32(width), int32(height), int32(w.screen.Stride), 1) // 1 = RGBA format
	w.pool.Destroy()

	// Ask compositor to create surface
	w.compositor.CreateSurface(w.surface)
	w.shell.GetShellSurface(w.shellsurface, w.surface)
	w.shellsurface.SetToplevel()

	// Make shell surface respond to pings
	pings := make(chan interface{})
	w.shellsurface.AddPingListener(pings)
	go func() {
		for p := range pings {
			ping := p.(gowl.ShellSurfacePing)
			w.shellsurface.Pong(ping.Serial)
		}
	}()

	go handleEvents(w)

	// Iterate
	go func () {
		for {
			w.display.Iterate()
		}
	} ()

	return
}

func (w *Window) SetTitle(title string) {
	w.shellsurface.SetTitle(title)
}

func (w *Window) SetSize(width, height int) {
}

func (w *Window) Size() (int, int) {
	return w.screen.Rect.Dx(), w.screen.Rect.Dy()
}

func (w *Window) Show() {
}

func (w *Window) Screen() draw.Image {
	return w.screen
}

func (w *Window) FlushImage(bounds ...image.Rectangle) {
	w.surface.Attach(w.buffer, 0, 0)

	for _, b := range bounds {
		w.surface.Damage(int32(b.Min.X), int32(b.Min.Y), int32(b.Dx()), int32(b.Dy()))
	}

	// Wait for redraw to finish
	cb := gowl.NewCallback()
	done := make(chan interface{})
	cb.AddDoneListener(done)
	w.surface.Frame(cb)
	func() {
		for {
			select {
			case <-done:
				return
			default:
				w.display.Iterate()
			}
		}
	}()
}

//func (w *Window) EventChan() <-chan interface{} {
//	return w.eventchan
//}

func (w *Window) Close() error {
	return nil
}

func waitForSync(display *gowl.Display) {
	cb := gowl.NewCallback()
	done := make(chan interface{})
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

func handleEvents(w *Window) {
	enterchan := make(chan interface{})
	leavechan := make(chan interface{})
	motionchan := make(chan interface{})
	buttonchan := make(chan interface{})

	w.pointer.AddEnterListener(enterchan)
	w.pointer.AddLeaveListener(leavechan)
	w.pointer.AddMotionListener(motionchan)
	w.pointer.AddButtonListener(buttonchan)

	for {
		select {
		case e := <-enterchan:
			enter := e.(gowl.PointerEnter)
			w.eventchan <- enter
		case _ = <-leavechan:
			w.eventchan <- 1
		case m := <-motionchan:
			motion := m.(gowl.PointerMotion)
//			w.eventchan <- 1
			w.eventchan <- motion
			fmt.Println("SENT",motion)
		case b := <-buttonchan:
			button := b.(gowl.PointerButton)
//			w.eventchan <- 1
			w.eventchan <- button
			fmt.Println("SENT",button)
		}
	}
}

func main() {
	w,_ := NewWindow(600,400)
	w.FlushImage(image.Rect(0,0,600,400))

	for e := range w.eventchan {
		fmt.Println("REC",e)
	}
}
