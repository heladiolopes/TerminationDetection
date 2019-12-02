package scholten

import (
  "log"
  "math/rand"
  "time"
)

func (scholten *Scholten) doWork(){
  scholten.currentState.Set(active)

  scholten.works += 1
  sleepTime := rand.Intn(maxSleepTime - minSleepTime) + minSleepTime
  log.Println("[ACTIVE] Going to work for ", sleepTime, "milliseconds")
  time.Sleep(time.Duration(sleepTime)*time.Millisecond)

  log.Println("[ACTIVE] Finished working! Now will wake up processes")

  for peerIndex := range scholten.peers {
    dice := rand.Float64()
    if dice < wakeUpProb && peerIndex != scholten.me{
      ok := scholten.sendBasic(peerIndex)
      if !ok {
        log.Println("Failed to send basic message to ", peerIndex)
      } else if peerIndex != scholten.me {
        scholten.addChild(peerIndex)
      }
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
