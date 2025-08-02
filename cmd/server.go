package main

import (
	"context"
	"log"
	"net/http"

	"simple_gin_server/configs"
	"simple_gin_server/internal/auth"
	"simple_gin_server/internal/match"
	"simple_gin_server/internal/profile"
	"simple_gin_server/internal/users"
	"simple_gin_server/pkg/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	routes     *Routes
	config     *configs.Config
	db         db.PgRepoInterface
}

// Конструктор для сервера
func NewServer(ctx context.Context) *Server {
	conf := configs.LoadConfig()

	// создаём экземпляр пула соединений на базе конфига и контекста
	db_pg, err := db.NewPgRepo(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}

	// создаём экземпляр reddis, используя config
	redisRepo := db.NewRedisRepo(ctx, conf)

	// Инициализация слоёв приложения

	//слой авторизации auth
	userRepository := users.NewUserRepository(db_pg)
	authService := auth.NewAuthService(userRepository, redisRepo)
	authHandler := auth.NewAuthHandler(authService, conf)

	//слой продукции match
	matchRepository := match.NewMatchRepository(db_pg)
	matchService := match.NewMatchService(matchRepository)
	matchHandler := match.NewMatchHandler(matchService, conf)

	//слой заказов profile
	profileRepository := profile.NewProfileRepository(db_pg)
	profileService := profile.NewProfileService(profileRepository, userRepository)
	ordersHandler := profile.NewProfileHandler(profileService)

	// создаём экземпляр роутера
	router := gin.Default()

	return &Server{
		router: router,
		routes: &Routes{
			Auth:    authHandler,
			Match:   matchHandler,
			Profile: ordersHandler,
			Config:  conf,
		},
		config: conf,
		db:     db_pg,
	}
}

// Метод для маршрутизации сервера
func (s *Server) SetUpRoutes() {
	s.router.SetTrustedProxies(nil)

	// Добавляем middleware для проброса контекста
	s.router.Use(func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "request_id", c.GetHeader("X-Request-ID"))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	s.routes.Setup(s.router)
}

// Метод для запуска сервера
func (s *Server) Run() error {
	s.SetUpRoutes()

	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}
	log.Println("Server is running on port 8080")
	return s.httpServer.ListenAndServe()
}

// Метод для graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	// Закрываем соединение с БД
	s.db.Close()

	// Останавливаем HTTP сервер
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server shutdown completed")
	return nil
}
