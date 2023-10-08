package main

import (
	"net"

	"github.com/modaniru/tages_test/gen/pkg"
	"github.com/modaniru/tages_test/internal/server"
	"github.com/modaniru/tages_test/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main(){
	imageService := service.NewImageService()
	imageServer := server.NewImageServiceServer(imageService)
	requestLimiter := server.NewRequestLimiter(imageServer)
	listener, _ := net.Listen("tcp", ":8080")
	s := grpc.NewServer()
	reflection.Register(s)
	pkg.RegisterImageServiceServer(s, requestLimiter)
	s.Serve(listener)
}