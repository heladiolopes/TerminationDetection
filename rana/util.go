package rana

import (
	"math/rand"
	"time"
)

// Returns a new value for timeout
func (rana *Rana) terminationTimeout() time.Duration {
	timeout := minActiveTimeoutMilli + rand.Intn(maxActiveTimeoutMilli-minActiveTimeoutMilli)
	return time.Duration(timeout) * time.Millisecond
}

// Reset the timer with a new timeout value
func (rana *Rana) resetTerminationTimeout() {
    rana.terminationTick = time.NewTimer(rana.terminationTimeout()).C
}
