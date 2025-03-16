package orchestrator

import (
	"log"
	"net/http"
)

func Run(port string, srvc *Service) {
	mux := http.NewServeMux()
	handlers := NewCustomHandlers(srvc)

	mux.Handle("/api/v1/calculate", LoggingMiddleware(http.HandlerFunc(handlers.ExpressionHandler)))
	mux.Handle("/api/v1/expressions", LoggingMiddleware(http.HandlerFunc(handlers.GetExpressions)))
	mux.Handle("/api/v1/expressions/", LoggingMiddleware(http.HandlerFunc(handlers.OneExpressionHandler)))
	mux.Handle("/internal/task", LoggingMiddleware(http.HandlerFunc(handlers.TaskHandler)))

	log.Println("Listening on " + port)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
