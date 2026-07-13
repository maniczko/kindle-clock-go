package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/y-yu/kindle-clock-go/inject"
	"github.com/y-yu/kindle-clock-go/presenter"
	"log/slog"
	"net/http"
	"os"
)

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					// スタックトレースを含むエラーを生成
					stackTrace := errors.WithStack(e)
					slog.Error("Recovered", "stacktrace", stackTrace)
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()
	r := newRouter(inject.ClockHandler(ctx), inject.HealthHandler(ctx))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := "0.0.0.0:" + port
	slog.Info("Kindle clock server started", "address", address)
	if err := http.ListenAndServe(address, r); err != nil {
		panic(err)
	}
}

func newRouter(clockHandler *presenter.ClockHandler, healthHandler *presenter.HealthHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(errorMiddleware)

	r.Get("/health", healthHandler.Handle)
	r.Get("/", clockHandler.Handle)
	r.Get("/clock", clockHandler.Handle)
	return r
}
