package main

import (
	"github.com/gin-gonic/gin"
	. "./v1"
)

func main() {
	gCTFRoute := gin.Default()
	v1:=gCTFRoute.Group("v1")
	{
		v1.POST("user",UserInfo)

	}

	gCTFRoute.Run(":8080")
}
