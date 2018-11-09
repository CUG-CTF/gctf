package main

import (
	. "./v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	gCTFRoute := gin.Default()
	gCTFRoute.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	v1 := gCTFRoute.Group("/v1")
	{

		userRoute := v1.Group("/user")
		{
			userRoute.POST("/info", UserInfo)
			userRoute.POST("/logout", Logout)
			userRoute.POST("/get_problem_list", GetProblemList)
			userRoute.POST("/get_score", GetScore)
			userRoute.POST("/start_problem/:name", StartProblem)
		}

		v1.GET("/get_users_rank", GetUsersRank)
		v1.POST("/get_teams_rank", GetTeamsRank)

	}

	gCTFRoute.Run(":8081")
}
