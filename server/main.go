// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/kiririmode/sandbox-grpc/greeter"
	"google.golang.org/grpc"
)

const (
	port            = ":50050"
	mutexSampleRate = 10
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// global lock
var mu sync.Mutex

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	now := ptypes.TimestampNow()
	log.Printf("SayHello: Hello %s at %s", req.Name, ptypes.TimestampString(now))
	return &greeter.HelloReply{
		Timestamp: now,
		Message:   "Hello " + req.Name,
	}, nil
}

// SayHellos implements helloworld.GreeterServer
func (s *server) SayHellos(req *greeter.HelloRequest, stream greeter.Greeter_SayHellosServer) error {
	done := make(chan interface{})
	ticker := time.NewTicker(time.Second * 60 * 60)
	defer ticker.Stop()
	time.AfterFunc(time.Second * 60 * 60, func() { close(done) })

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			/*
			if err := stream.Send(&greeter.HelloReply{
				Timestamp: ptypes.TimestampNow(),
				Message:   "Hello " + req.Name,
			}); err != nil {
				log.Fatalf("could not greet: %v", err)
				return err
			}
			*/
		}
	}
}

// AcuareLock acquires an global lock, and doesn't release it.
func (s *server) AcquireLock(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	log.Println("acquiring an lock")
	mu.Lock()
	// defer mu.Unlock()
	log.Println("acquired.")
	return &empty.Empty{}, nil
}

func main() {
	// profiling
	runtime.SetBlockProfileRate(mutexSampleRate)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
