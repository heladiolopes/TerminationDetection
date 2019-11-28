package scholten

// AppendEntryArgs is invoked by leader to replicate log entries; also used as
// heartbeat.
// Term  	- leader’s term
// leaderId - so follower can redirect clients
type ControlArgs struct {
    Sender   int
}


// AppendEntry is called by other instances of Scholten. It'll write the args received
// in the appendEntryChan.
func (rpc *RPC) Control(args *ControlArgs) error {
	rpc.scholten.appendEntryChan <- args
	return nil
}

// broadcastAppendEntries will send AppendEntry to all peers
func (scholten *Scholten) broadcastControl() {
    // TODO: Enviar control para todos

    // args := &AppendEntryArgs{
	// 	Term:     scholten.currentTerm,
	// 	LeaderID: scholten.me,
	// }
    //
	// for peerIndex := range scholten.peers {
	// 	if peerIndex != scholten.me {
	// 		go func(peer int) {
	// 			reply := &AppendEntryReply{}
	// 			ok := scholten.sendAppendEntry(peer, args, reply)
	// 			if ok {
	// 				reply.peerIndex = peer
	// 				replyChan <- reply
	// 			}
	// 		}(peerIndex)
	// 	}
	// }
}

// sendAppendEntry will send AppendEntry to a peer
func (scholten *Scholten) sendControl(peerIndex int, args *ControlArgs) bool {
	err := scholten.CallHost(peerIndex, "Control", args)
	if err != nil {
		return false
	}
	return true
}
