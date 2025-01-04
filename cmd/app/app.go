package app

import (
	"gestia/internal/app/gestia/handlers"
	"gestia/internal/app/gestia/repositories"
	"gestia/internal/app/gestia/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	_ "gestia/docs"

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

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           500, // Maximum value not ignored by any of major browsers
	}))

	imageRepository := repositories.NewMinioImageRepository()
	imageUsecase := usecases.NewImageUsecase(imageRepository)
	mainHandler := handlers.NewRootHandler(*imageUsecase)
	// Настройка маршрутов
	r.Get("/", mainHandler.HelloHandler)

	r.Route("/v1/images", func(r chi.Router) {
		r.Get("/", mainHandler.DownloadImagesHandler)
		r.Post("/", mainHandler.UploadImageHandler)
		r.Get("/{id}", mainHandler.GetImageHandler)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return &App{
		router: r,
	}, nil
}

func (a *App) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, a.router)
}
