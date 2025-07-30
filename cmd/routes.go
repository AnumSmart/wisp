package main

import (
	"simple_gin_server/configs"
	"simple_gin_server/internal/auth"
	"simple_gin_server/internal/match"
	"simple_gin_server/internal/profile"
	"simple_gin_server/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Auth    *auth.AuthHandler
	Match   *match.MatchHandler
	Profile *profile.ProfileHandler
	Config  *configs.Config
}

func (r *Routes) Setup(router *gin.Engine) {
	// Public routes
	public := router.Group("/")
	{
		public.POST("/register", middleware.ValidateAuthMiddleware(&auth.RegisterRequest{}), r.Auth.RegisterHandler) // эндпоинт для регистрации нового пользователя
		public.POST("/login", middleware.ValidateAuthMiddleware(&auth.LoginRequest{}), r.Auth.LoginHandler)          // эндпоинт для логина зарегестрированного пользователя (в ответе выдаётся access и refresh токены)
		public.POST("auth/refresh", r.Auth.ProcessRefreshTokenHandler)                                               // получение нового access токена при истечении его времени жизни, если refresh токен валиден и не в черном списке
	}

	// Authenticated routes
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware(r.Config))
	{
		authGroup.GET("/health", r.Auth.Check)          // health check, ручка-проверка, что все работатет
		authGroup.GET("/list", r.Auth.ListHandler)      // выводит список всех Email зарегестрированных юзеров
		authGroup.POST("/logout", r.Auth.LogoutHandler) // эндпоинт для logout (помощение refresh токена в черный список redis, удалени из БД)

		// User routes
		userGroup := authGroup.Group("/")
		userGroup.Use(middleware.RoleCheckMiddleware("user"))
		{
			r.setupUserRoutes(userGroup)
		}

		// Admin routes
		adminGroup := authGroup.Group("/")
		adminGroup.Use(middleware.RoleCheckMiddleware("admin"))
		{
			r.setupAdminRoutes(adminGroup)
		}
	}

}

func (r *Routes) setupUserRoutes(group *gin.RouterGroup) {
	group.POST("/get_new_access", r.Auth.ProcessRefreshTokenHandler)   // получение нового access токена при предоставлении валидного refresh токена в body
	group.POST("/profiles", r.Profile.CreateNewProfileHandler)         // создание нового профиля после авторизации (входящие данные JSON)
	group.GET("/profiles/me", r.Profile.GetMyProfileHandler)           // получение своего профиля(ответ в виде JSON)
	group.PATCH("/profiles/me", r.Profile.UpdateMyProfileHandler)      // обновление своего профиля
	group.DELETE("/profiles/me", r.Profile.DeleteMyProfileHandler)     // удаление своего профиля
	group.POST("/matches/search", r.Match.SearchMatchesHandler)        // получение списка совпадений по заданным критериям (входные данные JSON)
	group.POST("/matches/{id}/actions", r.Match.RegisterActionHandler) // регистрация действия пользователя (лайк/скип/жалоба)
	group.GET("/matches", r.Match.GetAcceptedMatchesHandler)           // получаем список совпадений, где 2-я сторона приняла запрос
	group.DELETE("/matches/{id}", r.Match.DeleteMetchByIdHandler)      // удалить совпадение по ID

}

func (r *Routes) setupAdminRoutes(group *gin.RouterGroup) {
	group.GET("/users", r.Auth.ListHandler)               // получить список  Email всех юзеров в базе
	group.DELETE("/users/{id}", r.Auth.DeleteUserHandler) // удалить конкретного юзера по id
}
