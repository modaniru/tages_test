package server

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"sync"

	"github.com/modaniru/tages_test/gen/pkg"
)

var ErrRequestLimited error = errors.New("request limited")

type RequestLimiter struct {
	//Для совместимости
	pkg.UnimplementedImageServiceServer
	imageServiceServer pkg.ImageServiceServer
	sayHello           func(func() error) error
	saveImage          func(func() error) error
	saveImageStream    func(func() error) error
	getImages          func(func() error) error
	getImagesStream    func(func() error) error
}

func NewRequestLimiter(server pkg.ImageServiceServer) *RequestLimiter {
	return &RequestLimiter{
		imageServiceServer: server,
		sayHello:           restrictions(10),
		saveImage:          restrictions(10),
		saveImageStream:    restrictions(10),
		getImages:          restrictions(100),
		getImagesStream:    restrictions(100),
	}
}

func (i *RequestLimiter) SayHello(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error) {
	var result *pkg.HelloReply
	err := i.sayHello(
		func() error {
			res, err := i.imageServiceServer.SayHello(ctx, request)
			result = res
			return err
		},
	)
	if err != nil {
		if !errors.Is(err, ErrRequestLimited) {
			log.Error("SayHello", log.String("status", fmt.Sprintf("error: %s", err.Error())))
		}
		return nil, err
	}
	log.Info("SayHello", log.String("status", "ok"))
	return result, nil
}

func (i *RequestLimiter) LoadImage(ctx context.Context, request *pkg.ImageRequest) (*pkg.Empty, error) {
	var result *pkg.Empty
	err := i.saveImage(
		func() error {
			res, err := i.imageServiceServer.LoadImage(ctx, request)
			result = res
			return err
		},
	)
	if err != nil {
		if !errors.Is(err, ErrRequestLimited) {
			log.Error("LoadImage", log.String("status", fmt.Sprintf("error: %s", err.Error())))
		}
		return nil, err
	}
	log.Info("LoadImage", log.String("status", "ok"))
	return result, nil
}

func (i *RequestLimiter) LoadImageStream(stream pkg.ImageService_LoadImageStreamServer) error {
	err := i.saveImageStream(
		func() error {
			return i.imageServiceServer.LoadImageStream(stream)
		},
	)
	if err != nil {
		if !errors.Is(err, ErrRequestLimited) {
			log.Error("LoadImageStream", log.String("status", fmt.Sprintf("error: %s", err.Error())))
		}
		return err
	}
	log.Info("LoadImageStream", log.String("status", "ok"))
	return nil
}

func (i *RequestLimiter) GetImages(ctx context.Context, empty *pkg.Empty) (*pkg.Images, error) {
	var result *pkg.Images
	err := i.getImages(
		func() error {
			res, err := i.imageServiceServer.GetImages(ctx, empty)
			result = res
			return err
		},
	)
	if err != nil {
		if !errors.Is(err, ErrRequestLimited) {
			log.Error("GetImages", log.String("status", fmt.Sprintf("error: %s", err.Error())))
		}
		return nil, err
	}
	log.Info("GetImages", log.String("status", "ok"))
	return result, nil
}

func (i *RequestLimiter) GetImagesStream(request *pkg.Empty, res pkg.ImageService_GetImagesStreamServer) error {
	err := i.getImages(
		func() error {
			err := i.imageServiceServer.GetImagesStream(request, res)
			return err
		},
	)
	if err != nil {
		if !errors.Is(err, ErrRequestLimited) {
			log.Error("GetImages", log.String("status", fmt.Sprintf("error: %s", err.Error())))
		}
		return err
	}
	log.Info("GetImages", log.String("status", "ok"))
	return nil
}

func (i *RequestLimiter) GetImagesInfo(ctx context.Context, request *pkg.Empty) (*pkg.ImagesInfo, error) {
	return i.imageServiceServer.GetImagesInfo(ctx, request)
}

func restrictions(limit int) func(func() error) error {
	count := 0
	var mutex sync.RWMutex
	return func(f func() error) error {
		mutex.Lock()
		count++
		if count > limit {
			count--
			mutex.Unlock()
			log.Error("requests limit error")
			return ErrRequestLimited
		}
		mutex.Unlock()

		defer func() {
			mutex.Lock()
			count--
			mutex.Unlock()
		}()
		err := f()
		if err != nil {
			return err
		}

		return nil
	}
}
