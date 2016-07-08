package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

const (
	PREPOSTDISTANCEM    = 5.0
	POSTFINISHDISTANCEM = 5.0
	MAXSPEED            = 10.0
	MINSPEED            = 0.1
)

type Measure struct {
	output     chan Muxable
	pre        chan BarrierEvent
	post       chan BarrierEvent
	finish     chan BarrierEvent
	preRunners []Runner
	runners    []Runner
}
type Runner struct {
	PreTime             time.Time
	PostTime            time.Time
	PrePostSpeed        float64
	EstimatedDuration   time.Duration
	EstimatedFinishTime time.Time
	FinishTime          time.Time
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

	pre := make(chan BarrierEvent)
	post := make(chan BarrierEvent)
	finish := make(chan BarrierEvent)

	NewBarriers(pre, post, finish)
	go keyboard(pre, post, finish)

	return &Measure{output: comms, pre: pre, post: post, finish: finish}
}

func (m *Measure) Loop() {
	for {
		select {
		case pre := <-m.pre:
			fmt.Printf("pre %s\n", pre)
			m.preBarrier(pre.Time)
		case post := <-m.post:
			fmt.Printf("post %s\n", post)
			m.postBarrier(post.Time)

		case finish := <-m.finish:
			fmt.Printf("finish %s\n", finish)
			m.finishBarrier(finish.Time)

		}
	}
}

//var starters []MeasurementStarted = make([]MeasurementStarted, 0, 5)

//var currentStarter *time.Time
//started := time.Now()
//ended := time.Now()
//duration := ended.Sub(started)

func (m *Measure) preBarrier(t time.Time) {
	m.preRunners = append(m.preRunners, Runner{PreTime: t})
	fmt.Println(len(m.preRunners))
}

func (m *Measure) postBarrier(t time.Time) {
	if len(m.preRunners) < 1 {
		return
	}

	r := m.preRunners[0]
	m.preRunners = m.preRunners[1:]

	dur := t.Sub(r.PreTime)
	r.PrePostSpeed = PREPOSTDISTANCEM / dur.Seconds()
	r.EstimatedDuration = time.Duration(r.PrePostSpeed / POSTFINISHDISTANCEM)
	r.EstimatedFinishTime = r.PreTime.Add(r.EstimatedDuration)

	if r.PrePostSpeed > MAXSPEED {
		fmt.Print("Too fast for this runner: %f, max: %f\n", r.PrePostSpeed, MAXSPEED)
		return
	}

	if r.PrePostSpeed < MINSPEED {
		fmt.Print("Too slow though pre and post barriers: %f, min: %f\n", r.PrePostSpeed, MINSPEED)
		return
	}

	m.output <- &MeasurementStarted{
		Started: r.PostTime,
		Speed:   r.PrePostSpeed,
	}

	r.PostTime = t
	m.runners = append(m.runners, r)
}

func (m *Measure) finishBarrier(t time.Time) {
	if len(m.runners) < 1 {
		return
	}

	lowestDiff := math.Abs(m.runners[0].EstimatedFinishTime.Sub(t).Seconds())
	index := 0
	for i, _ := range m.runners {
		diff := math.Abs(m.runners[i].EstimatedFinishTime.Sub(t).Seconds())
		if diff < lowestDiff {
			lowestDiff = diff
			index = i
		}
	}

	m.runners[index].FinishTime = t

	m.output <- &MeasurementEnded{
		Started:          m.runners[index].PostTime,
		Ended:            m.runners[index].FinishTime,
		Duration:         m.runners[index].FinishTime.Sub(m.runners[index].PostTime),
		DurationReadable: fmt.Sprintf("%s", m.runners[index].FinishTime.Sub(m.runners[index].PostTime)),
		Speed:            POSTFINISHDISTANCEM / m.runners[index].FinishTime.Sub(m.runners[index].PostTime).Seconds(),
	}

	m.runners = append(m.runners[:index], m.runners[index+1:]...)
}
