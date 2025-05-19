package main

import (
	"context"
	"log"
	"time"

	_ "grpc-lb-example/client_custom_lb/weighted_round_robin_lb"
	pb "grpc-lb-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address    = "example:///my-custom-service:1234"
	clientName = "Colin"
)

func main() {
	log.Println("客户端: 使用加权轮询负载均衡器连接服务:", address)

	// 创建gRPC连接，指定使用加权轮询负载均衡策略
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"weighted_round_robin"}`),
	)
	if err != nil {
		log.Fatalf("客户端: 连接失败: %v", err)
	}
	defer conn.Close()
	log.Println("客户端: 连接成功!")

	// 创建gRPC客户端
	c := pb.NewEchoClient(conn)

	// 发送10次请求
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
