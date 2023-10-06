package main

import (
	"fmt"
	"net"

	"github.com/modaniru/tages_test/gen/pkg"
	"github.com/modaniru/tages_test/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main(){

	server := server.NewImageServiceServer()
	listener, _ := net.Listen("tcp", ":8080")

	fmt.Println("ping pong")
	s := grpc.NewServer()
	reflection.Register(s)
	pkg.RegisterImageServiceServer(s, server)
	s.Serve(listener)
}