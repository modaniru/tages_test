package service

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/modaniru/tages_test/gen/pkg"
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

func (i *ImageService) GetImagesInfo() (*pkg.ImagesInfo, error) {
	dir, err := os.ReadDir("images")
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}
	images := []*pkg.ImageInfo{}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		i, err := f.Info()
		if err != nil {
			// TODO mb continue?
			return nil, fmt.Errorf("get file info error: %w", err)
		}
		stat := i.Sys().(*syscall.Stat_t)
		birthTime := time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
		images = append(images, &pkg.ImageInfo{
			Name:     f.Name(),
			CreateAt: birthTime.Format("2006-01-02 15:04:05"),
			UpdateAt: i.ModTime().Format("2006-01-02 15:04:05"),
		})
	}
	return &pkg.ImagesInfo{Images: images}, nil
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
