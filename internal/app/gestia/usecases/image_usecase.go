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

func (i *ImageUsecase) DownloadImages(offset int) (models.Image, error) {
	return i.imageRepository.GetImages(offset)
}
