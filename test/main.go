package main

import (
	"gowl"
	"time"
)


func main() {
	d := gowl.NewDisplay()
	d.Iterate()

	<-time.After(10e9)
}
