package interfaces

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const SERVICENAME = "com.xl.interface"

type HellServiceInterface = interface {
	SayHello(request string, reply *string) error
}

func RegisterHelloService(svc HellServiceInterface) error {
	return rpc.RegisterName(SERVICENAME, svc)
}

type HelloServiceClient struct {
	*rpc.Client
}

var _ *HelloServiceClient = (*HelloServiceClient)(nil)

func DialHelloService(network, add string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, add)
	if nil != err {
		return nil, err
	}
	return &HelloServiceClient{c}, nil
}

func DialHelloServiceJSON(network, add string) (*HelloServiceClient, error) {
	c, err := net.Dial(network, add)
	if nil != err {
		return nil, err
	}
	client :=rpc.NewClientWithCodec(jsonrpc.NewClientCodec(c))
	return &HelloServiceClient{client}, nil
}

func (p *HelloServiceClient) SayHello(request string, reply *string) error {
	return p.Call(SERVICENAME+".SayHello", request, reply)
}
