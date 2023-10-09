package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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

	i := 0
	for i < 500 {
		go func() {
			count := 0
			st, err := client.GetImagesStream(context.Background(), &pkg.Empty{})
			if err != nil {
				log.Fatal(err.Error())
			}
			for {
				images, err := st.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				count += len(images.Images)
			}
			fmt.Printf("Всего картинок: %d\n", count)
		}()
		i++
	}
	// img, err := os.ReadFile("client/images/maxresdefault.jpg")
	// now := time.Now()
	// _, err = client.GetImages(context.Background(), &pkg.Empty{})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(time.Now().Sub(now).Seconds())
	// for i := 200; i < 2000; i++ {
	// 	j := i
	// 	go func() {
	// 		_, err := client.LoadImage(context.Background(), &pkg.ImageRequest{Data: img, Name: fmt.Sprintf("%d.jpg", j)})
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 		}
	// 	}()
	// 	time.Sleep(time.Millisecond * 50)
	// }
	// for i := 0; i < 1000; i++ {
	// 	go func() {
	// 		_, err := client.GetImages(context.Background(), &pkg.Empty{})
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 		}
	// 	}()
	// 	time.Sleep(time.Millisecond)
	// }
	time.Sleep(time.Second * 60)
}
