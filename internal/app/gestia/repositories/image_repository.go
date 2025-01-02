package repositories

import "gestia/internal/app/gestia/models"

type IImageRepository interface {
	GetImages(limit, offset int) ([]models.Image, error)
	AddImage(models.Image) error
	GetImageByID(id string) (models.Image, error)
}
