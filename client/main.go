package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modaniru/tages_test/gen/pkg"
	"google.golang.org/grpc"
)

func main(){
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil{
		log.Fatal(err.Error())
	}
	fmt.Println("hrer")
	
	client := pkg.NewGreeterClient(conn)
	for i := 0; i < 1000; i++{
		go func ()  {
			j := i
			resp, err := client.SayHello(context.Background(), &pkg.HelloRequest{
				Name: "Данил",
			})
			if err != nil{
				fmt.Printf("%d: err = %s\n", j, err.Error())
			} else {
				fmt.Printf("%d: res = %s\n", j, resp.Message)
			}
		}()
		time.Sleep(time.Millisecond * 20)
	}
	time.Sleep(time.Second * 60)
}