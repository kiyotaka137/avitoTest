package http

import (
	"github.com/gin-gonic/gin"

	"avitoTest/internal/auth"
)

func NewRouter(
	teamHandler *TeamHandler,
	userHandler *UserHandler,
	prHandler *PullRequestHandler,
	statsHandler *StatsHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	authGroup := r.Group("/")
	authGroup.Use(auth.AuthMiddleware())

	adminGroup := authGroup.Group("/")
	adminGroup.Use(auth.AdminMiddleware())

	teamAuth := authGroup.Group("/team")
	{
		teamAuth.GET("/get", teamHandler.GetTeam) //получить команду с участниками
	}

	teamAdmin := adminGroup.Group("/team")
	{
		teamAdmin.POST("/add", teamHandler.AddTeam) //создать/обновить команду с участниками
	}

	usersAuth := authGroup.Group("/users")
	{
		usersAuth.GET("/getReview", userHandler.GetReviewPRs) //получить пры где пользователь ревьювер
	}

	usersAdmin := adminGroup.Group("/users")
	{
		usersAdmin.POST("/setIsActive", userHandler.SetIsActive) //установить флаг активности пользователю
	}

	pr := authGroup.Group("/pullRequest")
	{
		pr.POST("/create", prHandler.Create)     // создать пр и автоматически назначить до 2 ревьюверов из команды автора
		pr.POST("/merge", prHandler.Merge)       //пометить пр как merged(идемпотетен)
		pr.POST("/reassign", prHandler.Reassign) //переназначить конкретного ревьювера на другого из команды
	}

	stats := authGroup.Group("/stats")
	{
		stats.GET("/assignments", statsHandler.AssignmentStats) //сколько раз каждый пользователь сейчас назначен ревьювером в каких-либо PR и сколько ревьюверов сейчас назначено на каждый PR
	}

	return r
}
