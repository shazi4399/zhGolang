package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "zhGolang/hello-server/proto"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	// 开启端口
	listen, _ := net.Listen("tcp", ":9090")
	// 创建grpc服务
	grpcServer := grpc.NewServer()
	// 在grpc服务端中去注册我们自己编写的服务
	pb.RegisterSayHelloServer(grpcServer, &server{})

	// 启动grpc服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}