package service

import (
	"fmt"
	"os"
)

type ImageService struct{

}

func NewImageService() *ImageService{
	return &ImageService{}
}

func (i *ImageService) SaveImage(data []byte, name string) error{
	err := os.WriteFile(fmt.Sprintf("images/%s", name), data, 0644)
	if err != nil{
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}