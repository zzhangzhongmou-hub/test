package router

import (
	"test/handler"
	"test/middleware"
	"test/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.StaticFile("/", "./index.html")
	r.StaticFile("/homework.html", "./homework.html")

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

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
		auth.DELETE("/user/account", handler.DeleteAccount)
		auth.GET("/homework", handler.GetHomeworkList)
		auth.GET("/homework/:id", handler.GetHomeworkDetail)

		auth.POST("/submission", handler.Submit)
		auth.GET("/submission/my", handler.GetMySubmissions)
		auth.GET("/submission/excellent", handler.GetExcellentSubmissions)

		admin := auth.Group("/")
		admin.Use(middleware.RoleAuth(models.RoleAdmin))
		{
			admin.POST("/homework", handler.CreateHomework)
			admin.PUT("/homework/:id", handler.UpdateHomework)
			admin.DELETE("/homework/:id", handler.DeleteHomework)

			admin.GET("/submission/homework/:homework_id", handler.GetSubmissionsByHomework) // 查看作业的所有提交
			admin.PUT("/submission/:id/review", handler.Review)                              // 批改作业
			admin.PUT("/submission/:id/excellent", handler.MarkExcellent)
			admin.POST("/submission/:id/aiReview", handler.AIReview)

			admin.POST("/exam", handler.CreateExam)
			admin.GET("/exam/reviews", handler.GetMyReviews)
			admin.POST("/exam/review", handler.SubmitReview)
		}
	}
	return r
}
