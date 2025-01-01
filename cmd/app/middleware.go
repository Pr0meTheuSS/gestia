package app

import (
	"bytes"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Обертка над http.ResponseWriter для логирования статуса и тела ответа
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// Создание обертки с начальным статусом 200 OK
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// Перехват WriteHeader для записи статуса
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Перехват Write для записи тела ответа
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b) // Копируем данные в буфер
	return lrw.ResponseWriter.Write(b)
}

// Middleware для логирования с использованием zap
func NewZapMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Оборачиваем ResponseWriter
			lrw := newLoggingResponseWriter(w)

			// Передаем управление следующему обработчику
			next.ServeHTTP(lrw, r)

			// Логирование запроса и ответа
			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", lrw.statusCode),
				// zap.String("response_body", lrw.body.String()),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}
