package server

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/modaniru/tages_test/gen/pkg"
	"github.com/modaniru/tages_test/internal/service"
)

var ErrRequestLimited error = errors.New("request limited")

type ImageServiceServer struct {
	pkg.ImageRequest
	sayHello  func(func() error) error
	saveImage func(func() error) error
}

func NewImageServiceServer() *ImageServiceServer {
	return &ImageServiceServer{
		sayHello:  restrictions(10),
		saveImage: restrictions(10),
	}
}

func (i *ImageServiceServer) SayHello(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error) {
	var result *pkg.HelloReply
	err := i.sayHello(
		func() error{
			time.Sleep(time.Second * 1)
			result = &pkg.HelloReply{Message: "hello " + request.Name}
			return nil
		},
	)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func (i *ImageServiceServer) LoadImage(ctx context.Context, request *pkg.ImageRequest) (*pkg.Status, error) {
	fmt.Println("test")
	var result *pkg.Status
	err := i.saveImage(
		func() error{
			result = &pkg.Status{}
			service := service.NewImageService()
			err := service.SaveImage(request.Data, request.Name)
			if err != nil{
				return err
			}
			return nil
		},
	)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func restrictions(limit int) func(func() error) error {
	count := 0
	var mutex sync.RWMutex
	return func(f func() error) error{
		mutex.Lock()
		count++
		if count >= limit {
			count--
			mutex.Unlock()
			return ErrRequestLimited
		}
		mutex.Unlock()

		defer func(){
			mutex.Lock()
			count--
			mutex.Unlock()
		}()

		err := f()
		if err != nil{
			return err
		}

		return nil
	}
}
