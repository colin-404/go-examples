package main

import (
	"context"
	"log"
	"time"

	pb "grpc-custom-name-resolver/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// 重要：导入自定义解析器包，这会执行其 init() 函数来注册解析器
	_ "grpc-custom-name-resolver/resolver"
)

const (
	// 使用自定义 scheme "example" 和我们解析器中定义的服务名
	// 解析器会将 my-custom-service:1234 解析为 ["127.0.0.1:50051", "127.0.0.1:50052"]
	address    = "example:///my-custom-service:1234"
	clientName = "Colin"
)

func main() {
	log.Println("Client: Dialing with custom resolver address:", address)

	// Set up a connection to the server.
	// gRPC 会自动选择注册到 "example" scheme 的解析器
	// 默认情况下，gRPC 会在解析出的多个地址之间进行轮询 (round_robin) 负载均衡
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 如果你的解析器返回了 service config, 你可能不需要在这里显式设置
		// grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatalf("Client: Did not connect: %v", err)
	}
	defer conn.Close()
	log.Println("Client: Connected successfully!")

	c := pb.NewGreeterClient(conn)

	for range 5 { // 发送几次请求，观察负载均衡（如果服务端监听在不同端口）
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: clientName})
		if err != nil {
			log.Fatalf("Client: Could not greet: %v", err)
		}
		log.Printf("Client: Greeting from server: %s", r.GetMessage())
		time.Sleep(500 * time.Millisecond)
	}
}
