package server

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/modaniru/tages_test/gen/pkg"
	"github.com/modaniru/tages_test/internal/service"
)

// TODO remove simulator()
type ImageServiceServer struct {
	//Для совместимости
	pkg.UnimplementedImageServiceServer
	imageService *service.ImageService
}

func NewImageServiceServer(imageService *service.ImageService) *ImageServiceServer {
	return &ImageServiceServer{
		imageService: imageService,
	}
}

func (i *ImageServiceServer) SayHello(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error) {
	return &pkg.HelloReply{
		Message: "Hello, " + request.Name + "!",
	}, nil

}

func (i *ImageServiceServer) LoadImage(ctx context.Context, request *pkg.ImageRequest) (*pkg.Empty, error) {
	op := "internal.server.server.ImageServiceServer.LoadImage"
	err := i.imageService.SaveImage(request.Data, request.Name)
	simulate()
	if err != nil {
		return nil, fmt.Errorf("%s load image error: %w", op, err)
	}
	return &pkg.Empty{}, nil
}

func (i *ImageServiceServer) GetImages(context.Context, *pkg.Empty) (*pkg.Images, error) {
	op := "internal.server.server.ImageServiceServer.GetImages"
	images, err := i.imageService.GetImages()
	if err != nil {
		return nil, fmt.Errorf("%s get images error: %w", op, err)
	}
	result := pkg.Images{}
	for _, i := range images {
		result.Images = append(result.Images, &pkg.Image{
			Data: i.Data,
			Name: i.Name,
			Date: i.Date,
		})
	}
	return &result, nil
}

func simulate() {
	time.Sleep(time.Millisecond * time.Duration(1000+rand.Float64()*1000))
}
