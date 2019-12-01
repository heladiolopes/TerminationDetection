package scholten

import (
  "log"
)

func (scholten *Scholten) removeChild(child int) {
	log.Println("Removing child ", child)
	copy(scholten.children[i:], scholten.children[i+1:])
	scholten.children = scholten.children[:len(scholten.children)-1]
	scholten.ccp -= 1
}

func (scholten *Scholten) leaveTree(){
  if scholten.currentState == passive{
    if scholten.root {
  		log.Println("[PASSIVE] Termination detected!")
  		log.Println("[PASSIVE] Sending termination messages to all processes.")
  		scholten.broadcastFinish()
  	} else {
  		log.Println("[PASSIVE] Tree leaving condition met!")
  		log.Println("[PASSIVE] Sending control message to my father.")
  		scholten.sendControl(scholten.dad)
  	}
  }
}
