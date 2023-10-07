package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modaniru/tages_test/gen/pkg"
	"google.golang.org/grpc"
)

func main(){
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil{
		log.Fatal(err.Error())
	}
	fmt.Println("hrer")
	img, err := os.ReadFile("client/images/maxresdefault.jpg")
	if err != nil{
		log.Fatal(err.Error())
	}
	// data, name, err := image.Decode(bytes.NewReader(img))
	// if err != nil{
	// 	log.Fatal(err.Error())
	// }
	// fmt.Println("start load image: ", name)
	client := pkg.NewImageServiceClient(conn)
	// buf := new(bytes.Buffer)
	// png.Encode(buf, data)

	status, err := client.LoadImage(context.Background(), &pkg.ImageRequest{Data: img, Name: "test.jpg"})
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println(status)
	// for i := 0; i < 1000; i++{
	// 	go func ()  {
	// 		j := i
	// 		resp, err := client.SayHello(context.Background(), &pkg.HelloRequest{
	// 			Name: "Данил",
	// 		})
	// 		if err != nil{
	// 			fmt.Printf("%d: err = %s\n", j, err.Error())
	// 		} else {
	// 			fmt.Printf("%d: res = %s\n", j, resp.Message)
	// 		}
	// 	}()
	// 	time.Sleep(time.Millisecond * 20)
	// }
	// time.Sleep(time.Second * 60)
}