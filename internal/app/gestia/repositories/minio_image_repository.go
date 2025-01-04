package repositories

import (
	"bytes"
	"fmt"
	"gestia/internal/app/gestia/models"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

type MinioImageRepository struct {
	minioClient *minio.Client
	bucketName  string
	images      map[string]models.Image
	mu          sync.RWMutex
}

var (
	_ IImageRepository = (*MinioImageRepository)(nil)
)

func NewMinioImageRepository() IImageRepository {
	repo := MinioImageRepository{
		minioClient: &minio.Client{},
		bucketName:  "photos",
		images:      map[string]models.Image{},
		mu:          sync.RWMutex{},
	}

	var err error
	repo.minioClient, err = minio.New("localhost:9000", "admin", "admin123", false)
	if err != nil {
		log.Fatalf("MinIO connection error: %v", err)
	}

	// Проверка или создание бакета
	exists, err := repo.minioClient.BucketExists(repo.bucketName)
	if err != nil {
		return nil
	}
	if !exists {
		err = repo.minioClient.MakeBucket(repo.bucketName, "data")
		if err != nil {
			return nil
		}
	}

	return &repo
}

func (m *MinioImageRepository) AddImage(image models.Image) error {
	image.ID = uuid.NewString()
	_, err := m.minioClient.PutObject(m.bucketName, image.ID, bytes.NewReader(image.Data), int64(len(image.Data)), minio.PutObjectOptions{
		ContentType: "image/jpg", // Можно динамически определять тип
	})
	if err != nil {
		return err
	}

	image.Path = fmt.Sprintf("http://%s:%d/%s", "localhost", 9000, image.ID)

	m.mu.Lock()
	m.images[image.ID] = image
	m.mu.Unlock()

	return nil
}

func (m *MinioImageRepository) GetImageByID(id string) (models.Image, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.images[id], nil
}

func (m *MinioImageRepository) GetImages(limit int, offset int) ([]models.Image, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var images []models.Image
	for _, img := range m.images {
		images = append(images, img)
	}

	begin := min(offset, len(images))
	end := min(offset+limit, len(images))

	return images[begin:end], nil
}
