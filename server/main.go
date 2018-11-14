package main

import (
<<<<<<< HEAD
	. "./v1"
=======
	"./v1"
>>>>>>> upstream/master
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
	v1.ConfigRoute(gCTFRoute.Group("/v1"))
	gCTFRoute.Run(":8081")
}
