package repositories

import "gestia/internal/app/gestia/models"

type IImageRepository interface {
	GetImages(offset int) (models.Image, error)
	AddImage(models.Image) error
}
