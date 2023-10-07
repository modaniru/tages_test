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
	//Для совместимости
	pkg.UnimplementedImageServiceServer
	sayHello  func(func() error) error
	saveImage func(func() error) error
	getImages func(func() error) error
}

func NewImageServiceServer() *ImageServiceServer {
	return &ImageServiceServer{
		sayHello:  restrictions(10),
		saveImage: restrictions(10),
		getImages: restrictions(100),
	}
}

func (i *ImageServiceServer) SayHello(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error) {
	var result *pkg.HelloReply
	err := i.sayHello(
		func() error{
			service := service.NewImageService()
			fmt.Println(service.GetImages())
			result = &pkg.HelloReply{Message: "hello " + request.Name}
			return nil
		},
	)
	if err != nil{
		return nil, err
	}
	return result, nil
}

func (i *ImageServiceServer) LoadImage(ctx context.Context, request *pkg.ImageRequest) (*pkg.Empty, error) {
	var result *pkg.Empty
	err := i.saveImage(
		func() error{
			result = &pkg.Empty{}
			time.Sleep(time.Second)
			fmt.Println("test")
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

func (i *ImageServiceServer) GetImages(context.Context, *pkg.Empty) (*pkg.Images, error){
	var result *pkg.Images
	err := i.getImages(
		func() error{
			
			service := service.NewImageService()
			res, err := service.GetImages()
			if err != nil{
				return err
			}
			array := make([]*pkg.Image, 0, len(res))
			for _, i := range res{
				array = append(array, &pkg.Image{
					Data: i.Data,
					Name: i.Name,
					Date: i.Date,
				})
			}
			result = &pkg.Images{Images: array}
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
