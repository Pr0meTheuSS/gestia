package app

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestGetRoot(t *testing.T) {
	logger, _ := zap.NewProduction()
	app, _ := NewApp(logger)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %v; got %v", http.StatusOK, rec.Code)
	}

	expectedBody := "Hello, world!"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %v; got %v", expectedBody, rec.Body.String())
	}
}

func TestPostEmptyImage(t *testing.T) {
	logger, _ := zap.NewProduction()
	app, _ := NewApp(logger)

	body := bytes.NewBuffer([]byte{})
	req := httptest.NewRequest(http.MethodPost, "/v1/images/", body)
	req.Header.Add("Content-Type", "multipart/form-data")
	rec := httptest.NewRecorder()

	app.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %v; got %v", http.StatusBadRequest, rec.Code)
	}
}

func TestPostValidPngImage(t *testing.T) {
	// Создаем файл для загрузки
	filePath := "../../assets/test/images/valid_png_image.png" // Убедитесь, что файл существует
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	// Создаем тело запроса с файлом
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "valid_image.png")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file content: %v", err)
	}
	writer.Close()

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", "/v1/images/", body)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Устанавливаем заголовок Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Создаем тестовый сервер
	logger, _ := zap.NewProduction()
	app, _ := NewApp(logger)
	recorder := httptest.NewRecorder()

	// Выполняем запрос
	app.router.ServeHTTP(recorder, req)

	// Проверяем статус ответа
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
