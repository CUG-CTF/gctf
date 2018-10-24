package main

import "github.com/gin-gonic/gin"

func main() {
	gctfRoute:=gin.Default()

	gctfRoute.Run(":8080")
}