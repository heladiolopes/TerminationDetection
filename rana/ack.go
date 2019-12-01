package rana

// RequestVoteArgs is the struct that hold data passed to
// RequestVote RPC calls.
type AckArgs struct {
	Clock  int
	Sender int
}

// RequestVote is called by other instances of Raft. It'll write the args received
// in the requestVoteChan.
func (rpc *RPC) Ack(args *AckArgs, reply *AckArgs) error {
	rpc.rana.ackChan <- args
	return nil
}

// sendRequestVote will send RequestVote to a peer
func (rana *Rana) sendAck(peerIndex int, args *AckArgs) bool {
	reply := &AckArgs{}
	err := rana.CallHost(peerIndex, "Ack", args, reply)
	if err != nil {
		return false
	}
	return true
}
