package handlers

import (
	"encoding/json"
	"fmt"
	"gestia/internal/app/gestia/models"
	"gestia/internal/app/gestia/usecases"
	"io"
	"log"
	"net/http"
	"strconv"
)

type IRootHandler interface {
	HelloHandler(w http.ResponseWriter, r *http.Request)
	UploadImageHandler(w http.ResponseWriter, r *http.Request)
	DownloadImagesHandler(w http.ResponseWriter, r *http.Request)
	GetImageHandler(w http.ResponseWriter, r *http.Request)
}

type RootHandler struct {
	imageUsecase usecases.ImageUsecase
}

var (
	_ IRootHandler = (*RootHandler)(nil)
)

func NewRootHandler(imageUsecase usecases.ImageUsecase) IRootHandler {
	return &RootHandler{
		imageUsecase: imageUsecase,
	}
}

func (rh *RootHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

// @title File Download API
// @version 1.0
// @description API for downloading images.
// @BasePath /
// @Summary Download images
// @Description Download images file (JPEG/PNG).
// @Param limit query int true "Limit of images to fetch"
// @Param offset query int true	"Offset of images to fetch"
// @Produce text/plain
// @Success 200
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/images/ [get]

type ImageDTO struct {
	UUID string `json:"uuid"`
	URL  string `json:"url"`
}

func (rh *RootHandler) DownloadImagesHandler(w http.ResponseWriter, r *http.Request) {
	var limit int = 10
	var offset int
	var err error

	if limitString := r.URL.Query().Get("limit"); limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Download failed with parse limit %s.", err.Error()), http.StatusBadRequest)
			return
		}
	}

	if offsetString := r.URL.Query().Get("offset"); offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Download failed with parse offset %s.", err.Error()), http.StatusBadRequest)
			return
		}
	}

	images, err := rh.imageUsecase.DownloadImages(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Download failed with fetch images from repository: %s.", err.Error()), http.StatusBadRequest)
		return
	}

	imagesResponse := []ImageDTO{}
	for _, img := range images {
		imageDTO := ImageDTO{
			UUID: img.ID,
			URL:  fmt.Sprintf("http://localhost:9090/v1/images/%s", img.ID),
		}
		imagesResponse = append(imagesResponse, imageDTO)
	}
	responseJson, err := json.Marshal(imagesResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Download failed with marshal json: %s.", err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

// Новый обработчик для получения изображения по UUID
func (rh *RootHandler) GetImageHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// Отладочный вывод
	fmt.Println("All route parameters:", r.URL)
	fmt.Println(id)
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	img, err := rh.imageUsecase.GetImageByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch image: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Устанавливаем Content-Type в зависимости от типа изображения
	w.Header().Set("Content-Type", http.DetectContentType(img.Data))
	w.Write(img.Data)
}

// @title File Upload API
// @version 1.0
// @description API for uploading files.
// @BasePath /

// @Summary Upload an image
// @Description Upload an image file (JPEG/PNG).
// @Accept  multipart/form-data
// @Produce text/plain
// @Param file formData file true "File to upload"
// @Success 200 {string} string "File uploaded successfully!"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/images/ [post]
func (rh *RootHandler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {

	// Ограничение на размер загружаемого файла (например, 10 MB)
	const MaxUploadSize = 10 * 1024 * 1024

	// Проверяем размер загружаемого файла
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, fmt.Sprintf("The uploaded failed %s. Please upload a file less than 10MB.", err.Error()), http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve the file from the form.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Логирование информации о файле
	log.Printf("Uploaded File: %s\n", handler.Filename)
	log.Printf("File Size: %d\n", handler.Size)
	log.Printf("MIME Header: %v\n", handler.Header)

	newImage := models.Image{
		Name: handler.Filename,
		Data: []byte{},
	}

	// Чтение содержимого файла
	newImage.Data, err = io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the file.", http.StatusInternalServerError)
		return
	}

	// Проверка типа файла (например, изображение)
	fileType := http.DetectContentType(newImage.Data)
	if fileType != "image/jpeg" && fileType != "image/png" {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image.", http.StatusBadRequest)
		return
	}

	if err := rh.imageUsecase.UploadImage(newImage); err != nil {
		http.Error(w, fmt.Sprintf("Failed on usecase.UploadImage(). Error: %s", err.Error()), http.StatusBadRequest)
		return

	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully!"))
}
