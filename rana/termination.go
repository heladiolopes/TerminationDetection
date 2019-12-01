package rana

// AppendEntryArgs is invoked by leader to replicate log entries; also used as
// heartbeat.
// Term  	- leader’s term
// leaderId - so follower can redirect clients
type TerminationArgs struct {
    Initiator   int

}


// AppendEntry is called by other instances of Rana. It'll write the args received
// in the appendEntryChan.
func (rpc *RPC) Termination(args *TerminationArgs, reply *TerminationArgs) error {
    // args.Signature = make([]bool, len(rpc.rana.peers) + 1)
	rpc.rana.termChan <- args
	return nil
}

// broadcastAppendEntries will send AppendEntry to all peers
func (rana *Rana) broadcastTermination() {

	args := &TerminationArgs{
		Initiator:	rana.me,
	}

	for peerIndex := range rana.peers {
		// if peerIndex != rana.me {
			go func(peer int) {
				rana.sendTermination(peer, args)
			}(peerIndex)
		// }
	}
}

// sendAppendEntry will send AppendEntry to a peer
func (rana *Rana) sendTermination(peerIndex int, args *TerminationArgs) bool {
	reply := &TerminationArgs{}
	err := rana.CallHost(peerIndex, "Termination", args, reply)
	if err != nil {
		return false
	}
	return true
}
