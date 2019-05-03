// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	"github.com/kiririmode/sandbox-grpc/greeter"
	"google.golang.org/grpc"
)

const (
	port = ":50050"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	log.Printf("Received: %v", req.Name)
	return &greeter.HelloReply{Message: "Hello " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greeter.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
