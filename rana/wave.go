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
func (rpc *RPC) Wave(args *WaveArgs) error {
	rpc.rana.appendEntryChan <- args
    Signature = make(bool, len(rpc.rana.peers))

	return nil
}

// broadcastAppendEntries will send AppendEntry to all peers
func (rana *Rana) broadcastWave() {
    // TODO: Enviar wave para todos

    // args := &AppendEntryArgs{
	// 	Term:     rana.currentTerm,
	// 	LeaderID: rana.me,
	// }
    //
	// for peerIndex := range rana.peers {
	// 	if peerIndex != rana.me {
	// 		go func(peer int) {
	// 			reply := &AppendEntryReply{}
	// 			ok := rana.sendAppendEntry(peer, args, reply)
	// 			if ok {
	// 				reply.peerIndex = peer
	// 				replyChan <- reply
	// 			}
	// 		}(peerIndex)
	// 	}
	// }
}

// sendAppendEntry will send AppendEntry to a peer
func (rana *Rana) sendWave(peerIndex int, args *WaveArgs) bool {
	err := rana.CallHost(peerIndex, "Wave", args)
	if err != nil {
		return false
	}
	return true
}
