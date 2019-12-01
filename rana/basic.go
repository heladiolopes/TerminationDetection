package rana

import (
	"math/rand"
)

// RequestVoteArgs is the struct that hold data passed to
// RequestVote RPC calls.
type BasicArgs struct {
	Clock      int
	Sender     int

}

// RequestVote is called by other instances of Rana. It'll write the args received
// in the requestVoteChan.
func (rpc *RPC) Basic(args *BasicArgs) error {
	rpc.rana.basicChan <- args
	return nil
}

// broadcastRequestVote will send RequestVote to all peers
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
						rana.sendBasic(peer, wave)
					}(peerIndex)
				}
			}
		}

}

// sendRequestVote will send RequestVote to a peer
func (rana *Rana) sendBasic(peerIndex int, args *BasicArgs) bool {
	err := rana.CallHost(peerIndex, "Basic", args)
	if err != nil {
		return false
	}
	return true
}
