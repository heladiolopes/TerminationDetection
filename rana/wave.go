package rana

// WaveArgs is invoked by initiator to detect the termination,
// the others processes pass along the wave
// Clock  		- initiator timestap
// Initiator	- id from termination detection initiator
// Signature	- bool array with signature of processes that have already gone
	// through the wave
type WaveArgs struct {
	Clock       int
    Initiator   int

    Signature  [5]bool
}


// Wave is called by other instances of Rana. It'll write the args received
// in the waveChan. Requires response arguments, but not used.
func (rpc *RPC) Wave(args *WaveArgs, reply *WaveArgs) error {
	rpc.rana.waveChan <- args
	return nil
}

// broadcastWave will send Wave to all peers
func (rana *Rana) broadcastWave() {

	var signAux [5]bool
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

// sendWave will send Wave to a specific peer
func (rana *Rana) sendWave(peerIndex int, args *WaveArgs) bool {
	reply := &WaveArgs{}
	err := rana.CallHost(peerIndex, "Wave", args, reply)
	if err != nil {
		return false
	}
	return true
}
