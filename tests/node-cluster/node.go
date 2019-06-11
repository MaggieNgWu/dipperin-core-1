package node_cluster

import (
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/dipperin/dipperin-core/common"
	"fmt"
)

type Node struct {
	Client  *rpc.Client
	Address common.Address
}

// 新建一个rpc client
func newRpcClient(host string, port string) *rpc.Client {
	if host == "" {
		host = "127.0.0.1"
	}
	client, err := rpc.Dial(fmt.Sprintf("http://%v:%v", host, port))
	if err != nil {
		panic(err.Error())
	}
	return client
}