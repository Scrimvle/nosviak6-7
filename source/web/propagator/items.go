package propagator

import "sync"

var (
	// PropagatedFields displays everything which has been transmitted
	PropagatedFields map[string]any = make(map[string]any)
	mutex sync.RWMutex = sync.RWMutex{}
)