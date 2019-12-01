package scholten

import (
  "log",
  "math/rand",
  "time"
)

func (scholten *Scholten) doWork(){
  scholten.works += 1
  sleepTime = rand.Intn(maxSleepTime - minSleepTime) + minSleepTime
  log.Println("[ACTIVE] Going to work for ", sleepTime, "milliseconds")
  time.Sleep(sleepTime * time.Millisecond)

  ToWakeUp := make([]int, len(scholten.peers))
  order := rand.Perm(len(scholten.peers))
  dice := 1
  count := 0
  for dice := 1.0; dice > 0.5 || count < len(scholten.peers); dice = rand.Float64() {
    ToWakeUp = append(ToWakeUp, order[count])
    count += 1
  }
  log.Println("[ACTIVE] Finished working! Now will wake up ", len(ToWakeUp), "processes")
  for peerIndex := range ToWakeUp {
    ok := scholten.sendBasic(peerIndex)
    if !ok {
      log.Println("Failed to send basic message to ", peerIndex)
    }
  }

  scholten.works -= 1
  if scholten.works == 0 {
    scholten.currentState.Set(passive)
    if scholten.ccp == 0 {
      scholten.leaveTree()
    }
  }
}



// sendRequestVote will send RequestVote to a peer
func (scholten *Scholten) sendBasic(peerIndex int) bool {
	args := &BasicArgs{
		Sender: scholten.me
	}
	err := scholten.CallHost(peerIndex, "Basic", args)
	if err != nil {
		return false
	}
	return true
}
