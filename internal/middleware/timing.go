package middleware

import (
	"log"
	"net/http"
	"time"
)

func TimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		startTime := GetStartTime(r.Context())

		wrapped := &TimingResponseWriter{
			ResponseWriter: w,
			requestID:      requestID,
			startTime:      startTime,
		}

		next.ServeHTTP(wrapped, r)

		if !wrapped.headerWritten {
			duration := time.Since(startTime)
			w.Header().Set("X-Response-Time", duration.String())

			log.Printf("[%s] Request timing: %v", requestID, duration)
		}
	})
}

type TimingResponseWriter struct {
	http.ResponseWriter
	requestID     string
	startTime     time.Time
	headerWritten bool
}

func (trw *TimingResponseWriter) WriteHeader(code int) {
	if !trw.headerWritten {
		duration := time.Since(trw.startTime)
		trw.Header().Set("X-Response-Time", duration.String())

		log.Printf("[%s] Request timing: %v", trw.requestID, duration)
		trw.headerWritten = true
	}
	trw.ResponseWriter.WriteHeader(code)
}
