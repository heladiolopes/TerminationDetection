package rana

// AckArgs is invoked by the process to answer a Basic
// Clock  		- process timestap
// Sender		- id from message sender

type AckArgs struct {
	Clock  int
	Sender int

}

// Ack is called by other instances of Rana. It'll write the args received
// in the ackChan. Requires response arguments, but not used.
func (rpc *RPC) Ack(args *AckArgs, reply *AckArgs) error {
	rpc.rana.ackChan <- args
	return nil
}

// sendAck will send Ack to a specific peer
func (rana *Rana) sendAck(peerIndex int, args *AckArgs) bool {
	reply := &AckArgs{}
	err := rana.CallHost(peerIndex, "Ack", args, reply)
	if err != nil {
		return false
	}
	return true
}
