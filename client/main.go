package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/modaniru/tages_test/gen/pkg"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:80", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	client := pkg.NewImageServiceClient(conn)
	img, err := os.ReadFile("client/images/cat.jpg")
	now := time.Now()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(time.Now().Sub(now).Seconds())
	for i := 0; i < 2000; i++ {
		j := i
		go func() {
			_, err := client.LoadImage(context.Background(), &pkg.ImageRequest{Data: img, Name: fmt.Sprintf("%d.png", j)})
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
		time.Sleep(time.Millisecond * 50)
	}
	time.Sleep(time.Second * 60)
}
