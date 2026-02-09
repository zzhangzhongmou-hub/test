package router

import (
	"test/handler"
	"test/middleware"
	"test/models"

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
		auth.POST("/user/bindEmail", handler.BindEmail)
		auth.GET("/homework", handler.GetHomeworkList)
		auth.GET("/homework/:id", handler.GetHomeworkDetail)

		admin := auth.Group("/")
		admin.Use(middleware.RoleAuth(models.RoleAdmin))
		{
			admin.POST("/homework", handler.CreateHomework)
			admin.PUT("/homework/:id", handler.UpdateHomework)
			admin.DELETE("/homework/:id", handler.DeleteHomework)

			admin.POST("/exam", handler.CreateExam)
			admin.GET("/exam/reviews", handler.GetMyReviews)
			admin.POST("/exam/review", handler.SubmitReview)
		}
	}
	return r
}
