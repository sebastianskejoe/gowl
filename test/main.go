package main

import (
	"gowl"
	"time"
	"fmt"
)

type Display struct {
	display *gowl.Display
}

func (d *Display) globalListener(c chan interface{}) {
	for e := range c {
		glob := e.(gowl.DisplayGlobal)
		fmt.Println(glob.Name, glob.Iface, glob.Version)
	}
}

func main() {
	display := new(Display)
	display.display = gowl.NewDisplay()

	globchan := make(chan interface{})
	go display.globalListener(globchan)
	display.display.AddGlobalListener(globchan)

	display.display.Iterate()

	// Sync
	cb := gowl.NewCallback()
	done := make(chan interface{})
	cb.AddDoneListener(done)
	display.display.Sync(cb)
	for {
		select {
		case <-done:
			fmt.Println("Got done!")
			break
		default:
			display.display.Iterate()
		}
	}

	<-time.After(10e9)
}
