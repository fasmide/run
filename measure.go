package main

import (
	"fmt"
	"time"

	"encoding/json"
)

type Measure struct {
	output chan Muxable
}

type MeasurementStarted struct {
	Started time.Time `json:"started"`
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
}

func (m *MeasurementEnded) Marshal() *[]byte {
	encoded, _ := json.Marshal(m)
	return &encoded
}

func NewMeasure(comms chan Muxable) *Measure {
	return &Measure{output: comms}
}

func (m *Measure) Loop() {

	for {
		time.Sleep(1 * time.Second)

		started := time.Now()

		m.output <- &MeasurementStarted{Started: started}

		time.Sleep(5 * time.Second)

		ended := time.Now()
		duration := ended.Sub(started)

		m.output <- &MeasurementEnded{Started: started,
			Ended:            ended,
			Duration:         duration,
			DurationReadable: fmt.Sprintf("%s", duration),
		}
	}
}
