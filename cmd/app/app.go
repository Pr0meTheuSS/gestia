package app

import (
	handlers "gestia/internal/app/gestia"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	_ "gestia/docs" // Замените your_project_name на имя вашего модуля

	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	router *chi.Mux
}

func NewApp(logger *zap.Logger) (*App, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(NewZapMiddleware(logger))

	// Настройка маршрутов
	r.Get("/", handlers.RootHandler)
	r.Post("/v1/images/", handlers.UploadImageHandler)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return &App{
		router: r,
	}, nil
}

func (a *App) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, a.router)
}
