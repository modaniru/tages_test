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

func main(){
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil{
		log.Fatal(err.Error())
	}
	img, err := os.ReadFile("client/images/maxresdefault.jpg")
	if err != nil{
		log.Fatal(err.Error())
	}
	client := pkg.NewImageServiceClient(conn)

	
	for i := 0; i < 200; i++{
		j := i
		go func ()  {
			_, err := client.LoadImage(context.Background(), &pkg.ImageRequest{Data: img, Name: fmt.Sprintf("%d.jpg", j)})
			if err != nil{
				fmt.Println(err.Error())
			}
		}()
		time.Sleep(time.Millisecond * 50)
	}
	time.Sleep(time.Second * 60)
}