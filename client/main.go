// Package main implements a client for Greeter service.
package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/kiririmode/sandbox-grpc/greeter"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50050"
	defaultName = "world"
	timeout     = 20 * time.Second
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := greeter.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	stream, err := c.SayHellos(ctx, &greeter.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("SayHellos: %v", err)
		}
		log.Printf("Greeting: %s at %s", r.Message, ptypes.TimestampString(r.Timestamp))
	}
}
