package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
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

	// Чтение содержимого файла
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the file.", http.StatusInternalServerError)
		return
	}

	// Проверка типа файла (например, изображение)
	fileType := http.DetectContentType(fileBytes)
	if fileType != "image/jpeg" && fileType != "image/png" {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image.", http.StatusBadRequest)
		return
	}

	// Сохранение файла (опционально)
	// Создаем путь к файлу
	filePath := "../../assets/test/images/uploads/" + handler.Filename

	// Создаём директорию, если её нет
	if err := os.MkdirAll("../../assets/test/images/uploads/", os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Unable to create directory. Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Открываем файл, создавая его, если он отсутствует
	out, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to save the file. Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	fmt.Println(out.Name())
	// Записываем содержимое файла на диск
	if _, err := out.Write(fileBytes); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save the file. Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully!"))
}
