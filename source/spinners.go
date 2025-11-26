package source

import (
	"context"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
)

// SpinnerConfig represents information about each spinner
type SpinnerConfig struct {
	Frames []string `json:"frames"`
	Ticks  int      `json:"ticks"`
}

// SpinnerReceive is apart of the main spinner worker threads
type SpinnerReceive struct {
	ID        int
	Name      string
	LenFrames int
	FramePos  int
	Inherit   *SpinnerConfig
}

// Spinners stores all the receivers for the spinner workers
var Spinners []*SpinnerReceive = make([]*SpinnerReceive, 0)

// createSpinner will begin to start listening for spinner events.
func createSpinner(ctx context.Context, c []byte) error {
	var spinners map[string]*SpinnerConfig = make(map[string]*SpinnerConfig)
	if err := toml.Unmarshal(c, &spinners); err != nil {
		return err
	}

	Spinners = make([]*SpinnerReceive, 0)

	for name, spinner := range spinners {
		recv := &SpinnerReceive{
			Name: name,
			FramePos: 0,
			LenFrames: len(spinner.Frames),
			Inherit: spinner,
		}

		go func(receiver *SpinnerReceive) {
			Spinners = append(Spinners, receiver)
			recv.ID = len(Spinners) - 1

			ticker := time.NewTicker(time.Duration(receiver.Inherit.Ticks) * time.Millisecond)

			for {
				select {

				// context signals the worker to shutdown to allow us to allocate a new worker
				case <-ctx.Done():
					ticker.Stop()
					return

				case <-ticker.C:
					if recv.FramePos + 1 >= recv.LenFrames {
						Spinners[recv.ID].FramePos = 0
						continue
					}
	
					Spinners[recv.ID].FramePos++
				}
			}
		}(recv)
	}

	return nil
}

// GetSpinnerReceiver will look up the receiver with the given name
func GetSpinnerReceiver(name string) *SpinnerReceive {
	for _, spinner := range Spinners {
		if spinner.Name != name {
			continue
		}

		return spinner
	}

	return nil
}

// SmallestTickTime will find the smallest tick timer on a spinner
func SmallestTickTime() time.Duration {
	if len(Spinners) <= 0 {
		return 1000 * time.Millisecond
	}

	durations := make([]int, 0)
	for _, spinner := range Spinners {
		durations = append(durations, spinner.Inherit.Ticks)
	}

	sort.Ints(durations)
	return time.Duration(durations[0]) * time.Millisecond
}