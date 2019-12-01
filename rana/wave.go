package rana

// AppendEntryArgs is invoked by leader to replicate log entries; also used as
// heartbeat.
// Term  	- leader’s term
// leaderId - so follower can redirect clients
type WaveArgs struct {
	Clock       int
    Initiator   int

    // TODO: descobrir se vai funcionar na hora de testar
    Signature  []bool
}


// AppendEntry is called by other instances of Rana. It'll write the args received
// in the appendEntryChan.
func (rpc *RPC) Wave(args *WaveArgs, reply *WaveArgs) error {
    args.Signature = make([]bool, len(rpc.rana.peers) + 1)

	// for i := 1; i <= len(rpc.rana.peers); i++ {
	// 	arg.Signature[i] = false
	// }
	//
	// arg.Signature[args.Initiator] = true
	rpc.rana.waveChan <- args
	return nil
}

// broadcastAppendEntries will send AppendEntry to all peers
func (rana *Rana) broadcastWave() {

	signAux := make([]bool, len(rana.peers) + 1)
	for i := 1; i <= len(rana.peers); i++ {
		signAux[i] = false
	}
	signAux[rana.me] = true

	args := &WaveArgs{
		Clock:		rana.logicalClock,
		Initiator:	rana.me,

		Signature: signAux,
	}

	for peerIndex := range rana.peers {
		if peerIndex != rana.me {
			go func(peer int) {
				rana.sendWave(peer, args)
			}(peerIndex)
		}
	}
}

// sendAppendEntry will send AppendEntry to a peer
func (rana *Rana) sendWave(peerIndex int, args *WaveArgs) bool {
	reply := &WaveArgs{}
	err := rana.CallHost(peerIndex, "Wave", args, reply)
	if err != nil {
		return false
	}
	return true
}
