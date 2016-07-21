package main

import (
	"fmt"
	"time"

	"github.com/davecheney/gpio"
)

type Barrier struct {
	pin         gpio.Pin
	lastTrigger time.Time
}

func NewBarrier(pinIndex int, callback func(time.Time)) (*Barrier, error) {

	b := Barrier{}

	var err error

	b.pin, err = gpio.OpenPin(pinIndex, gpio.ModeInput)
	if err != nil {
		fmt.Printf("Error opening pin %d: %s\n", pinIndex, err.Error())
		return &b, err
	}

	err = b.pin.BeginWatch(gpio.EdgeFalling, func() {
		now := time.Now()
		if now.Sub(b.lastTrigger) > time.Millisecond*1000 {
			callback(now)
			b.lastTrigger = now
		}
	})

	if err != nil {
		fmt.Printf("Cannot watch pin %d: %s", pinIndex, err.Error())
		return &b, err
	}

	return &b, err
}
