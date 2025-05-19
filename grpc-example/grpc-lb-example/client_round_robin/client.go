package main

import (
	"context"
	"log"
	"time"

	pb "grpc-lb-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// 重要：导入自定义解析器包，这会执行其 init() 函数来注册解析器
	_ "grpc-lb-example/resolver"
)

const (
	address    = "example:///my-custom-service:1234"
	clientName = "Colin"
)

func main() {
	log.Println("Client: Dialing with custom resolver address:", address)

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatalf("Client: Did not connect: %v", err)
	}
	defer conn.Close()
	log.Println("Client: Connected successfully!")

	c := pb.NewEchoClient(conn)

	for range 10 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		r, err := c.Echo(ctx, &pb.EchoRequest{Message: clientName})
		if err != nil {
			log.Fatalf("Client: Could not greet: %v", err)
		}
		log.Printf("Client: Echo from server: %s", r.GetMessage())
		time.Sleep(500 * time.Millisecond)
	}
}
