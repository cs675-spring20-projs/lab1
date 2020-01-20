package serverless

import (
	"fmt"
	"net/rpc"
)

type WorkerRegisterArgs struct {
	WorkerAddr string
}

type ServiceRegisterArgs struct {
	ServiceName string
	ApiName     string
}

type RPCArgs struct {
	Name string
	Args []byte
}

func Call(srv string, rpcname string, args interface{}, reply interface{}) bool {
	c, errx := rpc.Dial("tcp", srv)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Printf("RPC call failed: %s\n", err)
	return false
}

// the following defines the generic plugin API
const (
	ServiceSymbolName = "Interface"
	PluginDir         = "../plugins"
)

type Interface interface {
	DoService(args []byte) error
}
