package scholten

import (
	"log"
)

// AppendEntryArgs is invoked by leader to replicate log entries; also used as
// heartbeat.
// Term  	- leader’s term
// leaderId - so follower can redirect clients
type ControlArgs struct {
	Sender int
}

type FinishArgs struct {
	Sender int
}

// AppendEntry is called by other instances of Scholten. It'll write the args received
// in the appendEntryChan.
func (rpc *RPC) Control(args *ControlArgs) error {
	rpc.scholten.controlChan <- args
	return nil
}

func (rpc *RPC) Finish(args *FinishArgs) error {
	rpc.scholten.finishChan <- args
	return nil
}

// broadcastAppendEntries will send AppendEntry to all peers
func (scholten *Scholten) broadcastFinish() {
	// TODO: Enviar finish para todos
	for peerIndex := range scholten.peers {
		if peerIndex != scholten.me {
      ok := scholten.sendFinish(peerIndex)
      if !ok {
        log.Println("Failed to send termination message to ", peerIndex)
      }
		}
	}
}

// sendAppendEntry will send AppendEntry to a peer
func (scholten *Scholten) sendControl(peerIndex int) bool {
	args := &ControlArgs{
		Sender: scholten.me,
	}
	err := scholten.CallHost(peerIndex, "Control", args)
	if err != nil {
		return false
	}
	return true
}

func (scholten *Scholten) sendFinish(peerIndex int) bool {
	args := &FinishArgs{
		Sender: scholten.me,
	}
	err := scholten.CallHost(peerIndex, "Finish", args)
	if err != nil {
		return false
	}
	return true
}
