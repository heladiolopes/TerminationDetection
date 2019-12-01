package scholten

import (
	"math/rand"
	"time"
)

// BroadcastInterval << ElectionTimeout << MTBF
// MTBF = Mean Time Between Failures

func (scholten *Scholten) terminationTimeout() time.Duration {
	timeout := minActiveTimeoutMilli + rand.Intn(maxActiveTimeoutMilli-minActiveTimeoutMilli)
	return time.Duration(timeout) * time.Millisecond
}

// func (raft *Raft) broadcastInterval() time.Duration {
// 	timeout := minElectionTimeoutMilli / 10
// 	return time.Duration(timeout) * time.Millisecond
// }

func (scholten *Scholten) resetTerminationTimeout() {
    scholten.terminationTick = time.NewTimer(scholten.terminationTimeout()).c
}
