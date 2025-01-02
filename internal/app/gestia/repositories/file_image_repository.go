package repositories

import (
	"gestia/internal/app/gestia/models"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type ImageRepository struct {
	images map[string]models.Image
	mu     sync.Mutex
}

var (
	_ IImageRepository = (*ImageRepository)(nil)
)

func NewImageRepository() IImageRepository {
	return &ImageRepository{
		images: map[string]models.Image{},
	}
}

func (i *ImageRepository) GetImages(limit, offset int) ([]models.Image, error) {
	var images []models.Image
	for _, img := range i.images {
		images = append(images, img)
	}
	offset = min(offset, len(images))

	return images[offset:min(offset+limit, len(images))], nil
}

func (i *ImageRepository) AddImage(image models.Image) error {
	filePath := "assets/test/images/uploads/" + image.Name
	var err error

	if err = os.MkdirAll("assets/test/images/uploads/", os.ModePerm); err != nil {
		return err
	}
	image.ID = uuid.NewString()
	image.Path, err = filepath.Abs(filePath)
	if err != nil {
		return err
	}

	i.mu.Lock()
	i.images[image.ID] = image
	i.mu.Unlock()

	out, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := out.Write(image.Data); err != nil {
		return err
	}

	return nil
}

// GetImageByID implements IImageRepository.
func (i *ImageRepository) GetImageByID(id string) (models.Image, error) {
	return i.images[id], nil
}
