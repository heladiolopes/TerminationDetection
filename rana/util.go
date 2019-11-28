package rana

import (
	"math/rand"
	"time"
)

// BroadcastInterval << ElectionTimeout << MTBF
// MTBF = Mean Time Between Failures

func (rana *Rana) terminationTimeout() time.Duration {
	timeout := minActiveTimeoutMilli + rand.Intn(maxActiveTimeoutMilli-minActiveTimeoutMilli)
	return time.Duration(timeout) * time.Millisecond
}

// func (raft *Raft) broadcastInterval() time.Duration {
// 	timeout := minElectionTimeoutMilli / 10
// 	return time.Duration(timeout) * time.Millisecond
// }

func (rana *Rana) resetTerminationTimeout() {
    rana.terminationTick = time.NewTimer(rana.terminationTimeout()).c
}
