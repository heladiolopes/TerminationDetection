package scholten

import (
	"log"
)

type TerminationArgs struct {
	Sender int
}

func (rpc *RPC) Termination(args *TerminationArgs, reply *Reply) error {
	rpc.scholten.terminationChan <- args
	return nil
}

func (scholten *Scholten) broadcastTermination() {
	for peerIndex := range scholten.peers {
		if peerIndex != scholten.me {
      ok := scholten.sendTermination(peerIndex)
      if !ok {
        log.Println("Failed to send termination message to ", peerIndex)
      }
		}
	}
}

func (scholten *Scholten) sendTermination(peerIndex int) bool {
	args := &TerminationArgs{
		Sender: scholten.me,
	}
	err := scholten.CallHost(peerIndex, "Termination", args)
	if err != nil {
		return false
	}
	return true
}
