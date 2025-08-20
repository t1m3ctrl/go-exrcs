package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware структура для middleware логирования
type LoggingMiddleware struct {
	handler http.Handler
}

// Logging создает middleware для логирования запросов
func Logging(next http.Handler) http.Handler {
	return &LoggingMiddleware{handler: next}
}

// ServeHTTP реализует интерфейс http.Handler
func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Создаем обертку для ResponseWriter чтобы захватить статус код
	wrapped := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Выполняем запрос
	l.handler.ServeHTTP(wrapped, r)

	// Логируем информацию о запросе
	duration := time.Since(start)
	log.Printf("[%s] %s %s - %d - %v - %s",
		start.Format("2006-01-02 15:04:05"),
		r.Method,
		r.URL.Path,
		wrapped.statusCode,
		duration,
		r.RemoteAddr,
	)
}

// responseWriter обертка для http.ResponseWriter для захвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader перехватывает статус код
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write записывает данные и устанавливает статус 200 если он не был установлен
func (w *responseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}
