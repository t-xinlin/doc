package main

import (
	"github.com/t-xinlin/doc/test/rpc/interfaces"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const SERVICENAME = "github.com/t-xinlin/rpc"

type HelloService struct {
}

func (helloSvice *HelloService) SayHello(request string, reply *string) error {
	log.Printf("Say " + request)
	*reply = "Say " + request
	return nil
}

func main() {
	//rpc.RegisterName("HelloService", new(HelloService))
	go rpcHttp()

	interfaces.RegisterHelloService(new(HelloService))

	listen, err := net.Listen("tcp", ":7777")
	if nil != err {
		log.Fatal("Listen error", err)
	}
	for {
		conn, err := listen.Accept()
		if nil != err {
			log.Fatal("Accept error", err)
		}

		//go rpc.ServeConn(conn)
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}

func rpcHttp() {
	interfaces.RegisterHelloService(new(HelloService))
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	http.ListenAndServe(":6666", nil)
}
