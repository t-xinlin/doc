package main

import (
	"context"
	"github.com/t-xinlin/doc/pkg/proto/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":50051"

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Send: %+v", "Hello"+in.Name)
	return &pb.HelloReply{Message: "Hello" + in.Name}, nil
}

func main() {
	li, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("fail to listen:%+v", port)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	err = s.Serve(li)
	if err != nil {
		log.Fatal("fail to listen:%+v", port)
	}
}
