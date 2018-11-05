package main

import (
	. "./v1"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
			userRoute.POST("/", UserInfo)
			userRoute.POST("/logout", Logout)
			userRoute.POST("/get_problem_list", GetProblemList)
			userRoute.POST("/get_score", GetScore)
			userRoute.POST("/start_problem/:name", StartProblem)
			//just del problem env
			userRoute.POST("submit_flag",SubmitFlag)
			userRoute.POST("del_problem",UserDelProblem)
		}
		adminRoute:=v1.Group("/admin")
		{
			//adminRoute.POST("/",adminInfo)
			//support add_user
			/*
			{
			[
			{"username":"xxx","passwd":"xxx","email":"xxx@xxx.com"},
			{"username":"xxx","passwd":"xxx","email":"xxx@xxx.com"}
			]
			}
			 */
			adminRoute.POST("/add_users",AddUsers)
			adminRoute.POST("/change_user_passwd",ChangeUserPasswd)
			//set user Problem one by one
			adminRoute.POST("/set_user_problem",SetUserProblem)
			adminRoute.POST("/random_all_users_problem",RandomAllUsersProblem)

			// train or test
			adminRoute.POST("/change_gctf_mode",ChangeGCTFMode)

			adminRoute.POST("/upload_problem",UploadProblem)
			adminRoute.POST("/change_problem",ChangeProblem)
			adminRoute.POST("/del_problem",DeleteProblem)
		}
		v1.GET("/get_users_rank", GetUsersRank)
		v1.POST("/get_teams_rank", GetTeamsRank)

	}
	gCTFRoute.Run(":8081")

}
