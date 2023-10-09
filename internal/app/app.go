package app

import (
	log "log/slog"
	"net"
	"os"

	"github.com/modaniru/tages_test/gen/pkg"
	"github.com/modaniru/tages_test/internal/server"
	"github.com/modaniru/tages_test/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func App() {
	configureLogger()
	log.Info("logger was configure")
	log.Info("DIP...")
	imageService := service.NewImageService()
	imageServer := server.NewImageServiceServer(imageService)
	requestLimiter := server.NewRequestLimiter(imageServer)
	listener, _ := net.Listen("tcp", ":8080")
	s := grpc.NewServer()
	reflection.Register(s)
	pkg.RegisterImageServiceServer(s, requestLimiter)
	log.Info("starting the server...")
	err := s.Serve(listener)
	if err != nil {
		log.Error("start server error")
		os.Exit(1)
	}
}
