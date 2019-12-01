package rana

import (
	"math/rand"
	"log"
)

// BasicArgs is invoked to activate other processes
// Clock  		- process timestap
// Sender		- id from message sender
type BasicArgs struct {
	Clock      int
	Sender     int

}

// Basic is called by other instances of Rana. It'll write the args received
// in the basicChan. Requires response arguments, but not used.
func (rpc *RPC) Basic(args *BasicArgs,  reply *BasicArgs) error {
	rpc.rana.basicChan <- args
	return nil
}

// broadcastBasic will send Wave to randon peers
func (rana *Rana) broadcastBasic() {

		args := &BasicArgs{
			Clock:		rana.logicalClock,
			Sender:		rana.me,
		}

		for peerIndex := range rana.peers {
			if peerIndex != rana.me {
				decision := rand.Float64()
				if decision < 0.4 {
					go func(peer int) {
						rana.sendBasic(peer, args)
						rana.ackToWait++
						log.Printf("[ACTIVE] Sending basic to %v.", peer)
					}(peerIndex)
				}
			}
		}

}

// sendBasic will send Basic to a specific peer
func (rana *Rana) sendBasic(peerIndex int, args *BasicArgs) bool {
	reply := &BasicArgs{}
	err := rana.CallHost(peerIndex, "Basic", args, reply)
	if err != nil {
		return false
	}
	return true
}
