package server

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/modaniru/tages_test/gen/pkg"
)

var ErrRequestLimited error = errors.New("request limited")

type RequestLimiter struct {
	//Для совместимости
	pkg.UnimplementedImageServiceServer
	imageServiceServer pkg.ImageServiceServer
	sayHello  func(func() error) error
	saveImage func(func() error) error
	getImages func(func() error) error
}

func NewRequestLimiter(server pkg.ImageServiceServer) *RequestLimiter{
	return &RequestLimiter{
		imageServiceServer: server,
		sayHello:  restrictions(10),
		saveImage: restrictions(10),
		getImages: restrictions(100),
	}
}

func (i *RequestLimiter) SayHello(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error) {
	var result *pkg.HelloReply
	err := i.sayHello(
		func() error{
			res, err := i.imageServiceServer.SayHello(ctx, request)
			result = res
			return err 
		},
	)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func (i *RequestLimiter) LoadImage(ctx context.Context, request *pkg.ImageRequest) (*pkg.Empty, error) {
	var result *pkg.Empty
	err := i.saveImage(
		func() error{
			res, err := i.imageServiceServer.LoadImage(ctx, request)
			result = res
			return err
		},
	)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func (i *RequestLimiter) GetImages(ctx context.Context, empty *pkg.Empty) (*pkg.Images, error){
	var result *pkg.Images
	err := i.getImages(
		func() error{
			res, err := i.imageServiceServer.GetImages(ctx, empty)
			result = res
			return err
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
		fmt.Println(count)
		err := f()
		if err != nil{
			return err
		}

		return nil
	}
}