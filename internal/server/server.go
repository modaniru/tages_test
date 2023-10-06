package server

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/modaniru/tages_test/gen/pkg"
)

var ErrRequestLimited error = errors.New("request limited")

type ImageServiceServer struct {
	pkg.ImageRequest
	sayHello  func(func()) error
	saveImage func(func()) error
}

// mustEmbedUnimplementedImageServiceServer implements pkg.ImageServiceServer.
func (*ImageServiceServer) mustEmbedUnimplementedImageServiceServer() {
	panic("unimplemented")
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
		func() {
			time.Sleep(time.Second * 1)
			result = &pkg.HelloReply{Message: "hello " + request.Name}
		},
	)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func (i *ImageServiceServer) LoadImage(ctx context.Context, request *pkg.ImageRequest) (*pkg.Status, error) {
	var result *pkg.Status
	err := i.sayHello(
		func() {
			result = &pkg.Status{}
			time.Sleep(1000)
		},
	)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func restrictions(limit int) func(func()) error {
	count := 0
	var mutex sync.RWMutex
	return func(f func()) error{
		mutex.Lock()
		count++
		if count >= limit {
			count--
			mutex.Unlock()
			return ErrRequestLimited
		}
		mutex.Unlock()

		f()

		mutex.Lock()
		count--
		mutex.Unlock()
		return nil
	}
}
