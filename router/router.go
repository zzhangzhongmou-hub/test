package router

import (
	"test/handler"
	"test/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", handler.Register)
		userGroup.POST("/login", handler.Login)
	}
	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/user/profile", handler.GetProfile)
	}
	return r
}
