package server

import (
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
}

func NewServer() Server {
	cfg := configs.New()
	db, _ := repositories.GetPool(cfg)

	return Server{router: chi.NewRouter(), cfg: cfg, db: db}
}

func (s *Server) setupRoutes() {
	userRepository := repositories.NewUserRepository(s.db)
	authService := services.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	roomRepo := repositories.NewRoomRepository(s.db)
	messageRepo := repositories.NewMessageRepository(s.db)
	roomServ := services.NewChatRoomService(roomRepo, messageRepo)
	roomHadler := handlers.NewChatRoomHandler(roomServ)

	ws := websocket.NewWebsocketHandler(roomServ)

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

func (s *Server) Start(port string) {
	s.setupMiddlewares()
	s.setupRoutes()

	slog.Info("Starting server on :" + port)
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		slog.Error("Server failed to start: " + err.Error())
	}
}
