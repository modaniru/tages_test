package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/modaniru/tages_test/gen/pkg"
)

type ImageServiceServer struct{
	pkg.GreeterServer
	f func(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error)
}

func NewImageServiceServer() *ImageServiceServer{
	return &ImageServiceServer{f: say()}
}

func (i *ImageServiceServer) SayHello(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error){
	return i.f(ctx, request)
}

func say() func(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error){
	count := 0
	var mutex sync.RWMutex
	return func(ctx context.Context, request *pkg.HelloRequest) (*pkg.HelloReply, error) {
		mutex.Lock()
		count++
		if count > 10{
			count--
			mutex.Unlock()
			return nil, errors.New("error")
		}
		mutex.Unlock()
		
		fmt.Println(count)
		time.Sleep(time.Millisecond * time.Duration((rand.Float64() * 500)))

		mutex.Lock()
		count--
		mutex.Unlock()
		
		return &pkg.HelloReply{Message: "hello " + request.Name}, nil
	}
}