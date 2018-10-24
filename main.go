package main

import (
	. "./v1"
	"github.com/gin-gonic/gin"
)

func main() {
	gCTFRoute := gin.Default()
	v1 := gCTFRoute.Group("/v1")
	{
		v1.POST("/user", UserInfo)
		v1.POST("/logout",Logout)
		v1.POST("/get_problem",GetProblem)
		v1.POST("/get_score",GetScore)

	}

	gCTFRoute.Run(":8080")
}
