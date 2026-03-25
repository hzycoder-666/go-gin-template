package routes

import (
	"github.com/gin-gonic/gin"
	"hzycoder.com/go-gin-template/internal/handler"
	"hzycoder.com/go-gin-template/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
	}

	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		api.GET("/user/me", handler.GetUserInfo)
		api.GET("/user/info/:name", handler.GetUser)
	}

	return r
}
