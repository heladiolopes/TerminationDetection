package rana

// TerminationArgs is called by the initiator to communicate termination
// Initiator	- id from termination detection initiator
type TerminationArgs struct {
    Initiator   int

}

// Termination is called by other instances of Rana. It'll write the args received
// in the termChan. Requires response arguments, but not used.
func (rpc *RPC) Termination(args *TerminationArgs, reply *TerminationArgs) error {
    // args.Signature = make([]bool, len(rpc.rana.peers) + 1)
	rpc.rana.termChan <- args
	return nil
}

// broadcastTermination will send Termination to all peers, myself included
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

// sendTermination will send Termination to a specifiv peer
func (rana *Rana) sendTermination(peerIndex int, args *TerminationArgs) bool {
	reply := &TerminationArgs{}
	err := rana.CallHost(peerIndex, "Termination", args, reply)
	if err != nil {
		return false
	}
	return true
}
