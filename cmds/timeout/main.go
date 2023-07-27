package main

import (
	_ "context"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(1 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("hi"))
	})

	r.Get("/long", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		processTime := time.Duration(rand.Intn(4)+1) * time.Second

		select {
		case <-ctx.Done():
			return

		case <-time.After(processTime):
			// The above channel simulates some hard work.
		}

		w.Write([]byte("done"))
	})

	restServer := &http.Server{
		Addr:    ":3333",
		Handler: r,
	}

	restServer.ListenAndServe()
}
