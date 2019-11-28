package rana

import "net/rpc"

// RPC is the struct that implements the methods exposed to RPC calls
type RPC struct {
	rana *Rana
}

// CallHost will communicate to another host through it's RPC public API.
func (rana *Rana) CallHost(index int, method string, args interface{}) error {
	var (
		err    error
		client *rpc.Client
	)

	client, err = rpc.Dial("tcp", rana.peers[index])
	if err != nil {
		return err
	}

	defer client.Close()

	err = client.Call("RPC."+method, args, nil)

	if err != nil {
		return err
	}

	return nil
}
