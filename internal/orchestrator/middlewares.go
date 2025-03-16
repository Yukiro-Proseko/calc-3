package orchestrator

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Ошибка при чтении тела запроса: %v", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rr, r)

		log.Printf(
			"Method: %s | Endpoint: %s | Body: %s | Status Code: %d",
			r.Method,
			r.URL.Path,
			string(bodyBytes),
			rr.statusCode,
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}
