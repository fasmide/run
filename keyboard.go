package main

import (
	"os"
	"os/exec"
	"time"
)

func keyboard(pre, post, finish chan BarrierEvent, m *Measure) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "sane").Run()

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)

		if b[0] == 49 {
			pre <- BarrierEvent{Time: time.Now()}

		}
		if b[0] == 50 {
			post <- BarrierEvent{Time: time.Now()}

		}
		if b[0] == 51 {

			finish <- BarrierEvent{Time: time.Now()}

		}
		if b[0] == 52 {
			m.Flush()

		}
	}

}
