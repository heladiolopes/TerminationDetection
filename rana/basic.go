package rana

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
    // TODO: definir lógica para envio de mensagens basic

    // args := &RequestVoteArgs{
	// 	CandidateID: rana.me,
	// 	Term:        rana.currentTerm,
	// }
    //
	// for peerIndex := range rana.peers {
	// 	if peerIndex != rana.me { // exclude self
	// 		go func(peer int) {
	// 			reply := &RequestVoteReply{}
	// 			ok := rana.sendRequestVote(peer, args, reply)
	// 			if ok {
	// 				reply.peerIndex = peer
	// 				replyChan <- reply
	// 			}
	// 		}(peerIndex)
	// 	}
	// }
}

// sendRequestVote will send RequestVote to a peer
func (rana *Rana) sendBasic(peerIndex int, args *BasicArgs) bool {
	err := rana.CallHost(peerIndex, "Basic", args)
	if err != nil {
		return false
	}
	return true
}
