package main

import (
	"fmt"
	"time"
)

type BarrierEvent struct {
	Time time.Time
}

func NewBarriers(pre, post, finish chan BarrierEvent) {

	NewBarrier(17, func(t time.Time) {
		fmt.Println("Vi skal arbejde med preBarrieren: %s", t)
		pre <- BarrierEvent{Time: t}
	})

	NewBarrier(27, func(t time.Time) {
		fmt.Println("Vi skal arbejde med postBarrieren: %s", t)
		post <- BarrierEvent{Time: t}

	})

	NewBarrier(22, func(t time.Time) {
		fmt.Println("Vi skal arbejde med finishBarrieren: %s", t)
		finish <- BarrierEvent{Time: t}

	})
}
