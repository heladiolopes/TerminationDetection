package scholten

// RequestVoteArgs is the struct that hold data passed to
// RequestVote RPC calls.
type BasicArgs struct {
	Sender     int
}

// RequestVote is called by other instances of Scholten. It'll write the args received
// in the requestVoteChan.
func (rpc *RPC) Basic(args *BasicArgs) error {
	rpc.scholten.basicChan <- args
	return nil
}

// broadcastRequestVote will send RequestVote to all peers
// func (scholten *Scholten) broadcastBasic() {
    // TODO: definir lógica para envio de mensagens basic

    // args := &RequestVoteArgs{
	// 	CandidateID: scholten.me,
	// 	Term:        scholten.currentTerm,
	// }
    //
	// for peerIndex := range scholten.peers {
	// 	if peerIndex != scholten.me { // exclude self
	// 		go func(peer int) {
	// 			reply := &RequestVoteReply{}
	// 			ok := scholten.sendRequestVote(peer, args, reply)
	// 			if ok {
	// 				reply.peerIndex = peer
	// 				replyChan <- reply
	// 			}
	// 		}(peerIndex)
	// 	}
	// }
// }

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
