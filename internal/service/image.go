package service

import (
	"fmt"
	"os"

	"github.com/djherbis/times"
	"github.com/modaniru/tages_test/gen/pkg"
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
		t, err := times.Stat(fmt.Sprintf("images/%s", i.Name()))
		if err != nil {
			return nil, err
		}
		bt := "system not support btime"
		if t.HasBirthTime() {
			bt = t.BirthTime().Format("2006-01-02 15:04:05")
		}
		images = append(images, &pkg.ImageInfo{
			Name:     f.Name(),
			CreateAt: bt,
			UpdateAt: i.ModTime().Format("2006-01-02 15:04:05"),
		})
	}
	return &pkg.ImagesInfo{Images: images}, nil
}

func (i *ImageService) GetImage(name string) ([]byte, error) {
	res, err := os.ReadFile(fmt.Sprintf("images/%s", name))
	if err != nil {
		return nil, fmt.Errorf("read image with name %s error: %w", name, err)
	}
	return res, nil
}
