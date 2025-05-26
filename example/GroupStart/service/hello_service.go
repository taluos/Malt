package service

import (
	"context"
	"log"

	pb "github.com/taluos/Malt/example/test_proto"
)

// GreeterServer 实现 Greeter 服务
type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello 实现 SayHello RPC 方法
func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("收到来自 %s 的问候请求", req.GetName())
	return &pb.HelloReply{Message: "你好, " + req.GetName()}, nil
}

// NewGreeterServer 创建一个新的 GreeterServer
func NewGreeterServer() *GreeterServer {
	return &GreeterServer{}
}
