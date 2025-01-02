package usecases

import (
	"gestia/internal/app/gestia/models"
	"gestia/internal/app/gestia/repositories"
)

type ImageUsecase struct {
	imageRepository repositories.IImageRepository
}

func NewImageUsecase(repo repositories.IImageRepository) *ImageUsecase {
	return &ImageUsecase{
		imageRepository: repo,
	}
}

func (i *ImageUsecase) UploadImage(image models.Image) error {
	return i.imageRepository.AddImage(image)
}

func (i *ImageUsecase) DownloadImages(limit, offset int) ([]models.Image, error) {
	return i.imageRepository.GetImages(limit, offset)
}

func (i *ImageUsecase) GetImageByID(id string) (models.Image, error) {
	return i.imageRepository.GetImageByID(id)
}
