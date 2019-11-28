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


// AppendEntry is called by other instances of Raft. It'll write the args received
// in the appendEntryChan.
func (rpc *RPC) Wave(args *WaveArgs) error {
	rpc.raft.appendEntryChan <- args
    Signature = make(bool, len(rpc.rana.peers))

	return nil
}

// broadcastAppendEntries will send AppendEntry to all peers
func (rana *Rana) broadcastWave() {
    // TODO: Enviar wave para todos

    // args := &AppendEntryArgs{
	// 	Term:     raft.currentTerm,
	// 	LeaderID: raft.me,
	// }
    //
	// for peerIndex := range raft.peers {
	// 	if peerIndex != raft.me {
	// 		go func(peer int) {
	// 			reply := &AppendEntryReply{}
	// 			ok := raft.sendAppendEntry(peer, args, reply)
	// 			if ok {
	// 				reply.peerIndex = peer
	// 				replyChan <- reply
	// 			}
	// 		}(peerIndex)
	// 	}
	// }
}

// sendAppendEntry will send AppendEntry to a peer
func (raft *Raft) sendWave(peerIndex int, args *WaveArgs) bool {
	err := raft.CallHost(peerIndex, "Wave", args)
	if err != nil {
		return false
	}
	return true
}
