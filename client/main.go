// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

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

	limit := make(chan struct{}, 10)
	var wg sync.WaitGroup
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func(i int) {
			limit <- struct{}{}

			log.Printf("acquireLock: %d", i)
			_, err := c.AcquireLock(context.Background(), &empty.Empty{})
			if err != nil {
				log.Fatalf("AcquireLock: %v", err)
			}
			wg.Done()
			<-limit
		}(i)
	}
	wg.Wait()
}
