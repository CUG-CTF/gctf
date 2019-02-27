package v1

import "github.com/gin-gonic/gin"

func ConfigRoute(r *gin.RouterGroup) {
	userRoute := r.Group("/user")
	{
		userRoute.Use(checkSessionMiddleware)

		userRoute.POST("/", UserInfo)
		userRoute.POST("/logout", Logout)
		userRoute.POST("/get_problem_list", GetProblemList)
		userRoute.POST("/get_score", GetScore)
		userRoute.POST("/start_problem/:id", StartProblem)
		//just del problem env
		userRoute.POST("submit_flag", SubmitFlag)
		userRoute.POST("del_problem", UserDelProblem)
	}
	adminRoute := r.Group("/admin")
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
		adminRoute.Use(checkSessionMiddleware)
		adminRoute.POST("/add_users", AddUsers)
		adminRoute.POST("/change_user_passwd", ChangeUserPasswd)
		//set user Problem one by one
		adminRoute.POST("/set_user_problem", SetUserProblem)
		adminRoute.POST("/random_all_users_problem", RandomAllUsersProblem)

		// train or test
		adminRoute.POST("/change_gctf_mode", ChangeGCTFMode)

		adminRoute.POST("/upload_problem", UploadProblem)
		adminRoute.POST("/change_problem", ChangeProblem)
		adminRoute.POST("/del_problem", DeleteProblem)
	}
	r.GET("/get_users_rank", GetUsersRank)
	r.POST("/get_teams_rank", GetTeamsRank)
	r.POST("/login",Login)
}
