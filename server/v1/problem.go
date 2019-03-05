package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartProblem(c *gin.Context) {
	type startProblem struct {
		Username   string `json:"username"`
		Token      string `json:"token"`
		Problem_ID int    `json:"problem_id"`
	}
	var sp startProblem
	err := c.BindJSON(&sp)
	if err != nil {
		log.Println("user/StartProblem: errot to bind json" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to start bind json"})
	}
	//TODO:insert DB start container
	//c,err:=gctfConfig.DockerClient.ContainerCreate(context.Background(),nil,nil,nil,"")
}

func GetProblemList(c *gin.Context) {

}
