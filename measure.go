package main

import (
	"fmt"

	"time"

	"encoding/json"
)

const (
	PREPOSTDISTANCEM    = 0.15
	POSTFINISHDISTANCEM = 0.5
	MAXSPEED            = 10.0
	MINSPEED            = 0.1
)

type Measure struct {
	output chan Muxable
}

type MeasurementStarted struct {
	Started time.Time `json:"started"`
	Speed   float64   `json:"speed"`
}

func (m *MeasurementStarted) Marshal() *[]byte {
	encoded, _ := json.Marshal(m)
	return &encoded

}

type MeasurementEnded struct {
	Id               string        `json:"id" bson:"-"`
	Started          time.Time     `json:"started"`
	Ended            time.Time     `json:"ended"`
	Duration         time.Duration `json:"durationNs"`
	DurationReadable string        `json:"durationHumanReadable"`
	Speed            float64       `json:"speed"`
}

func (m *MeasurementEnded) Marshal() *[]byte {
	encoded, _ := json.Marshal(m)
	return &encoded
}

func NewMeasure(comms chan Muxable) *Measure {
	return &Measure{output: comms}
}

func (m *Measure) Loop() {

	var starters []MeasurementStarted = make([]MeasurementStarted, 0, 5)

	var currentStarter *time.Time

	err = preStartPin.BeginWatch(gpio.EdgeFalling, func() {
		fmt.Printf("Hej\n", "hej")
		tmp := time.Now()
		currentStarter = &tmp
	})

	err = postStartPin.BeginWatch(gpio.EdgeFalling, func() {

	})

	err = finishPin.BeginWatch(gpio.EdgeFalling, func() {

	})
	// ended := time.Now()
	// duration := ended.Sub(started)

	// m.output <- &MeasurementEnded{Started: started,
	// 	Ended:            ended,
	// 	Duration:         duration,
	// 	DurationReadable: fmt.Sprintf("%s", duration),
	//  }
}
