package main

import (
	"fmt"

	"github.com/davecheney/gpio"
)

type Barriers struct {
	preBarrier    Barrier
	postBarrier   Barrier
	finishBarrier Barrier
}

func NewBarriers(awfw chan Runner) *Barriers {
	runners := make([]Runner, 0, 5)

	preBarrier = NewBarrier(17, func(t time.Time) {
		fmt.Printf("Vi skal arbejde med preBarrieren: %s", t)
		runners = append(runners, Runner{PreTime: t})
	})

	postBarrier = NewBarrier(27, func(t time.Time) {
		fmt.Printf("Vi skal arbejde med postBarrieren: %s", t)
		for runner, index := range runners {
			err := runner.postBarrier(t)
			if err != nil {
				// remove the runner from the slice
				// Soooo, we are changing an slice that we are ranging over? ..
				runners = append(runners[:index], runners[index+1:])
				continue
			}

			awfw <- runner
			break

		}
	})

	finishBarrier = NewBarrier(22, func(t time.Time) {
		fmt.Printf("Vi skal arbejde med finishBarrieren: %s", t)
		for runner, _ := range runners {
			runner.finishBarrier(t)
		}

	})
}

type Runner struct {
	PreTime      time.Time
	PostTime     time.Time
	PrePostSpeed float64
	FinishTime   time.Time
}

func (r *Runner) postBarrier(t time.Time) error {

	dur := r.PreTime.Sub(t)
	r.PrePostSpeed = PREPOSTDISTANCEM / dur.Seconds()

	if r.PrePostSpeed > MAXSPEED {
		// too fast
		return error.Error(fmt.Sprintf("Too fast for this runner: %f, max: %f", r.PrePostSpeed, MAXSPEED))
	}

	if r.PrePostSpeed < MINSPEED {
		return error.Error(fmt.Sprintf("Too slow though pre and post barriers: %f, min: %f\n", r.PrePostSpeed, MINSPEED))
	}

	r.PostTime = t

}

func (r *Runner) finishBarrier(t time.Time) error {

	minDuration := (r.PrePostSpeed - 3) / POSTFINISHDISTANCEM
	maxDuration := (r.PrePostSpeed + 3) / POSTFINISHDISTANCEM

	duration := t.Sub(r.PostTime)

	if duration.Seconds() < minDuration {
		return error.Error(fmt.Sprintf("This runner is way too fast: %s, %f min\n", duration, minDuration))
	}

	if duration.Seconds() > maxDuration {
		fmt.Printf("This runner took to long: %s, %f allowed\n", duration, maxDuration)
		return
	}

	m.output <- &MeasurementEnded{Started: mStarted.Started,
		Ended:            ended,
		Duration:         duration,
		DurationReadable: fmt.Sprintf("%s", duration),
	}

}
