package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code"`
	RequestID string `json:"request_id"`
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r.Context())

				log.Printf("[%s] PANIC: %v\n%s", requestID, err, debug.Stack())

				if rw, ok := w.(*ResponseWriter); ok && rw.StatusCode != 0 {
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				errorResp := ErrorResponse{
					Error:     "Internal server error occurred",
					Code:      "INTERNAL_ERROR",
					RequestID: requestID,
				}

				json.NewEncoder(w).Encode(errorResp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
