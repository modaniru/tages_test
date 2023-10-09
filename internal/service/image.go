package service

import (
	"fmt"
	"os"

	"github.com/modaniru/tages_test/internal/entity"
)

type ImageService struct{}

func NewImageService() *ImageService {
	return &ImageService{}
}

func (i *ImageService) SaveImage(data []byte, name string) error {
	err := os.WriteFile(fmt.Sprintf("images/%s", name), data, 0644)
	if err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}

func (i *ImageService) GetImages() ([]entity.Image, error) {
	dir, err := os.ReadDir("images")
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}
	images := []entity.Image{}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		i, err := f.Info()
		if err != nil {
			return nil, fmt.Errorf("get file info error: %w", err)
		}
		date := i.ModTime().Format("2006-01-02 15:04:05")
		data, err := os.ReadFile(fmt.Sprintf("images/%s", i.Name()))
		if err != nil {
			return nil, fmt.Errorf("read file error: %w", err)
		}
		images = append(images, entity.Image{Name: f.Name(), Date: date, Data: data})
	}
	return images, nil
}
