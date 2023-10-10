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
	saveImage          func(func() error) error
	saveImageStream    func(func() error) error
	getImageStream     func(func() error) error
	getImagesInfo      func(func() error) error
}

func NewRequestLimiter(server pkg.ImageServiceServer) *RequestLimiter {
	return &RequestLimiter{
		imageServiceServer: server,
		saveImage:          restrictions(10),
		saveImageStream:    restrictions(2),
		getImageStream:     restrictions(10),
		getImagesInfo:      restrictions(100),
	}
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

func (i *RequestLimiter) GetImageStream(req *pkg.GetImageRequest, stream pkg.ImageService_GetImageStreamServer) error {
	err := i.getImageStream(
		func() error {
			return i.imageServiceServer.GetImageStream(req, stream)
		},
	)
	if err != nil {
		if !errors.Is(err, ErrRequestLimited) {
			log.Error("GetImageStream", log.String("status", fmt.Sprintf("error: %s", err.Error())))
		}
		return err
	}

	log.Info("GetImageStream", log.String("status", "ok"))
	return nil
}

func (i *RequestLimiter) GetImagesInfo(ctx context.Context, request *pkg.Empty) (*pkg.ImagesInfo, error) {
	result := &pkg.ImagesInfo{}
	err := i.getImagesInfo(
		func() error {
			res, err := i.imageServiceServer.GetImagesInfo(ctx, request)
			result = res
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		if !errors.Is(err, ErrRequestLimited) {
			log.Error("GetImagesInfo", log.String("status", fmt.Sprintf("error: %s", err.Error())))
		}
		return nil, err
	}

	log.Info("GetImagesInfo", log.String("status", "ok"))
	return result, nil
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
