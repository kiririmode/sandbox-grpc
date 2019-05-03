// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/kiririmode/sandbox-grpc/greeter"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50050"
	defaultName = "world"
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

	done := make(chan interface{})
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	time.AfterFunc(10*time.Second, func() { close(done) })

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Println("Current time: ", t)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.SayHello(ctx, &greeter.HelloRequest{Name: name})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Greeting: %s", r.Message)
		}
	}
}
