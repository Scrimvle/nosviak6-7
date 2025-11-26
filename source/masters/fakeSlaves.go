package masters

import (
	"Nosviak4/source"
	"Nosviak4/source/masters/terminal"
	"math/rand"
	"time"
)

// fakeSlavesWorker implements the brand new fake slave generator
func fakeSlavesWorker() {
	for range time.NewTicker(time.Duration(source.OPTIONS.Ints("fake_slaves", "sleep")) * time.Millisecond).C {
		if !source.OPTIONS.Bool("fake_slaves", "enabled") {
			continue
		}

		terminal.Mutex.Lock()
		frequency := rand.Intn(source.OPTIONS.Ints("fake_slaves", "max_change") - source.OPTIONS.Ints("fake_slaves", "min_change")) + source.OPTIONS.Ints("fake_slaves", "max_change")
				
		switch rand.Intn(2) % 2 == 0 {

		// increase the fake count frequency
		case true: 
			if terminal.FakeSlaves + frequency > source.OPTIONS.Ints("fake_slaves", "maximum") {
				terminal.FakeSlaves -= frequency
			} else {
				terminal.FakeSlaves += frequency
			}

		case false:
			if terminal.FakeSlaves - frequency < source.OPTIONS.Ints("fake_slaves", "minimum") {
				terminal.FakeSlaves += frequency
			} else {
				terminal.FakeSlaves -= frequency
			}
		}

		terminal.Mutex.Unlock()
	}
}