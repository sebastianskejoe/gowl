package main

import (
	"gowl"
)


func main() {
	c := make(chan bool)
	d := gowl.NewDisplay(c)
	d.Iterate()
	d.Sync()

	d.Compositor.CreateSurface()

	<-c
}
