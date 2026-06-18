package middleware

import (
	"log"
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.Size += size
	return size, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		startTime := GetStartTime(r.Context())

		log.Printf("[%s] Started %s %s from %s", requestID, r.Method, r.URL.Path, r.RemoteAddr)

		wrapped := &ResponseWriter{
			ResponseWriter: w,
			StatusCode:     200,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(startTime)
		log.Printf("[%s] Completed %s %s - Status: %d, Size: %d bytes, Duration: %v",
			requestID, r.Method, r.URL.Path, wrapped.StatusCode, wrapped.Size, duration)
	})
}
