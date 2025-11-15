package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/i474232898/chatserver/configs"
	"github.com/i474232898/chatserver/internal/app/handlers"
	"github.com/i474232898/chatserver/internal/app/middlewares"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/services"
	"github.com/i474232898/chatserver/internal/app/websocket"
	"github.com/swaggest/swgui/v5emb"
	"gorm.io/gorm"
)

type Server struct {
	router *chi.Mux
	cfg    *configs.AppConfigs
	db     *gorm.DB
	server *http.Server
}

func NewServer() Server {
	cfg := configs.New()
	db, _ := repositories.GetPool(cfg)

	return Server{chi.NewRouter(), cfg, db, nil}
}

func (s *Server) setupRoutes() {
	userRepository := repositories.NewUserRepository(s.db)
	authService := services.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	healthcheckHandler := handlers.NewHealthcheckHandler()

	roomRepo := repositories.NewRoomRepository(s.db)
	messageRepo := repositories.NewMessageRepository(s.db)
	roomServ := services.NewChatRoomService(roomRepo, messageRepo)
	roomHadler := handlers.NewChatRoomHandler(roomServ)

	ws := websocket.NewWebsocketHandler(roomServ)

	s.router.Get("/healthcheck", healthcheckHandler.Healthcheck)
	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/signin", authHandler.Signin)
	})
	s.router.Route("/user", func(r chi.Router) {
		r.Use(middlewares.JWTAuthMiddleware([]byte("secret")))
		r.Get("/me", userHandler.Me)
	})
	s.router.Route("/rooms", func(r chi.Router) {
		r.Use(middlewares.JWTAuthMiddleware([]byte("secret")))
		r.Get("/", roomHadler.ListRooms)
		r.Post("/", roomHadler.CreateRoom)
		r.Post("/direct", roomHadler.DirectMessage)
	})
	// s.router.Get("/ws", websocket.WebsocketHandler)
	s.router.Route("/ws", func(r chi.Router) {
		// r.Use(middlewares.JWTAuthMiddleware([]byte("secret")))

		//ws/room/{roomID}?token=JWT
		r.Get("/room/{roomID}", ws.JoinChatRoomHandler)
	})

	s.router.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi.yaml")
	})
	s.router.Mount("/docs", v5emb.NewHandler(
		"Chat Server API Docs",
		"/openapi.yaml",
		"/docs",
	))
}

func (s *Server) setupMiddlewares() {
	r := s.router

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares.ContentTypeJSONMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

}

func (s *Server) Start(ctx context.Context, port string) {
	s.setupMiddlewares()
	s.setupRoutes()
	s.server = &http.Server{Addr: ":" + port, Handler: s.router}

	slog.Info("Starting server on :" + port)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server error: " + err.Error())
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		slog.Error("Server shutdown error: " + err.Error())
		return err
	}
	slog.Info("Server shutdown complete")
	return err
}
