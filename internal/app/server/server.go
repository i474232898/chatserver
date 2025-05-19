package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/i474232898/chatserver/internal/app/handlers"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/services"
	"github.com/i474232898/chatserver/pkg/db"
)

type Server struct {
	router *chi.Mux
}

func NewServer() Server {
	return Server{router: chi.NewRouter()}
}

func (s *Server) setupRoutes() {
	db := db.GetPool()

	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)

	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
	})
}

func (s *Server) setupMiddlewares() {
	r := s.router

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

}

func (s *Server) Start(port string) {
	s.setupMiddlewares()
	s.setupRoutes()

	slog.Info("Starting server on :" + port)
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		slog.Error("Server failed to start: " + err.Error())
	}
}
