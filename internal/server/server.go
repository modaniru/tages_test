package server

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/modaniru/tages_test/gen/pkg"
	"github.com/modaniru/tages_test/internal/service"
)

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
	// simulate()
	if err != nil {
		return nil, fmt.Errorf("%s load image error: %w", op, err)
	}

	return &pkg.Empty{}, nil
}

func (i *ImageServiceServer) LoadImageStream(stream pkg.ImageService_LoadImageStreamServer) error {
	op := "internal.server.server.ImageServiceServer.LoadImageStream"
	for {
		request, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&pkg.Empty{})
		}

		if err != nil {
			return fmt.Errorf("%s stream.Recv() error: %w", op, err)
		}

		err = i.imageService.SaveImage(request.Data, request.Name)
		if err != nil {
			return fmt.Errorf("%s load image error: %w", op, err)
		}
	}
}

func (i *ImageServiceServer) GetImageStream(req *pkg.GetImageRequest, stream pkg.ImageService_GetImageStreamServer) error {
	for _, name := range req.GetNames() {
		data, err := i.imageService.GetImage(name)
		if err != nil {
			return err
		}

		err = stream.Send(&pkg.GetImageResponse{Data: data})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *ImageServiceServer) GetImagesInfo(ctx context.Context, request *pkg.Empty) (*pkg.ImagesInfo, error) {
	return i.imageService.GetImagesInfo()
}
