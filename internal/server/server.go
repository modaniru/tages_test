package server

import (
	"context"

	"github.com/modaniru/tages_test/gen/pkg"
)

type ImageServiceServer struct{
	pkg.GreeterServer
}

func NewImageServiceServer() *ImageServiceServer{
	return &ImageServiceServer{}
}

func (i *ImageServiceServer) SayHello(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error){
	return &pkg.HelloReply{
		Message: "hello " + request.Name + "!",
	}, nil
}