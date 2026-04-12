package routes

import (
	"github.com/gin-gonic/gin"
	"hzycoder.com/go-gin-template/internal/handler"
	"hzycoder.com/go-gin-template/internal/middleware"
	"hzycoder.com/go-gin-template/internal/model"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.RecoveryWithBizError())
	r.Use(middleware.Logger())
	r.Use(middleware.BizErrorHandler())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/register", handler.Register)
		}

		v1 := api.Group("/v1")
		v1.Use(middleware.JWTAuth())
		{
			v1.GET("/user/me", middleware.RequireRole(model.RoleMember), handler.GetUserInfo)
			v1.GET("/user/info/:name", handler.GetUser)
			v1.POST("/generate/img", handler.GenerateImage)
			v1.POST("/generate/img/action", handler.UpdateImageAction)
			v1.GET("/generate/img", handler.QueryGeneratedImageById)
			v1.GET("/generate/img/list", handler.QueryGeneratedImageList)
			v1.POST("/generate/img/modal", handler.UpdateImageModal)
		}
	}

	return r
}
