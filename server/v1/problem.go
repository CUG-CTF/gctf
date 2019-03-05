package v1

import (
	"github.com/gin-gonic/gin"
)


func StartProblem( c *gin.Context){
	type startProblem struct {
		Username string `json:"username"`
		Token string `json:"token"`
	}
	var sp startProblem
	c.BindJSON(&sp)
}

func GetProblemList(c *gin.Context)  {

}