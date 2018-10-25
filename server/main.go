package main

import (
	. "./v1"
	"github.com/gin-gonic/gin"
)


func main() {
	gCTFRoute := gin.Default()
	v1 := gCTFRoute.Group("/v1")
	{

		userRoute:= v1.Group("/user")
		{
			userRoute.POST("/info", UserInfo)
			userRoute.POST("/logout", Logout)
			userRoute.POST("/get_problem_list", GetProblemList)
			userRoute.POST("/get_score", GetScore)
			userRoute.POST("/start_problem/:name", StartProblem)
		}

		v1.POST("/get_rank",GetRank)
	}



	gCTFRoute.Run(":8080")
}
