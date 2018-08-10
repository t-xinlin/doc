package main

import (
	"google.golang.org/grpc"
	"log"
	"os"
	"context"
	"github.com/t-xinlin/test/proto/pb"
)

const address = "localhost:50051"
const defaultName = "world"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if nil != err {
		log.Fatal("did not connect %+v", address)
	}
	if nil != conn {
		defer conn.Close()
	}
	c := pb.NewGreeterClient(conn)
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	for {
		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if nil != err {
			log.Fatal("could not greet:%+v", err)
		}
		log.Printf("Re: %s", r.Message)
	}

}
