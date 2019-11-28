package scholten

import "net/rpc"

// RPC is the struct that implements the methods exposed to RPC calls
type RPC struct {
	scholten *Scholten
}

// CallHost will communicate to another host through it's RPC public API.
func (scholten *Scholten) CallHost(index int, method string, args interface{}) error {
	var (
		err    error
		client *rpc.Client
	)

	client, err = rpc.Dial("tcp", scholten.peers[index])
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
