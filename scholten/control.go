package scholten

import (
	"log"
)

type ControlArgs struct {
	Sender int
}

func (rpc *RPC) Control(args *ControlArgs, reply *Reply) error {
	rpc.scholten.controlChan <- args
	return nil
}

func (scholten *Scholten) sendControl(peerIndex int) bool {
	args := &ControlArgs{
		Sender: scholten.me,
	}
	err := scholten.CallHost(peerIndex, "Control", args)
	if err != nil {
		return false
	}
	return true
}
