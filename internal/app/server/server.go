package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/i474232898/chatserver/internal/app/auth"
)

func Start(port string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/signup", auth.SignupH)

	log.Printf("starting server on :%s", port)
	http.ListenAndServe(":"+port, r)
}
