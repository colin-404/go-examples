package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "grpc-lb-example/proto"

	"google.golang.org/grpc"
)

var addrs = []string{"localhost:50051", "localhost:50052"}

type server struct {
	pb.UnimplementedEchoServer
	addr string
}

func (s *server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: fmt.Sprintf("%s (from %s)", req.Message, s.addr)}, nil
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	for _, addr := range addrs {
		go startServer(addr)
	}

	<-stop
	log.Println("Shutting down servers...")

}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening at %v", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{addr: addr})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
