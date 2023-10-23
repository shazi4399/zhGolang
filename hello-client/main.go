package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "zhGolang/hello-server/proto"
)

func main() {
	// 连接到server端，此处禁用安全传输，没有加密和验证
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 建立链接
	client := pb.NewSayHelloClient(conn)

	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhanghao ddddd"})
	fmt.Println(resp.GetMessage())
}
