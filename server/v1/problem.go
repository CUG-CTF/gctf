package v1

import (
	"../gctfConfig"
	"../model"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartProblem(c *gin.Context) {
	type startProblem struct {
		Username   string `json:"username"`
		Token      string `json:"token"`
		Problem_ID int64  `json:"problem_id"`
	}
	var sp startProblem
	err := c.BindJSON(&sp)
	if err != nil {
		log.Println("user/StartProblem: errot to bind json" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to start bind json"})
	}
	var up model.UserProblems
	var p model.Problems
	var u model.User
	u.Username = sp.Username
	p.Id = sp.Problem_ID
	h, err := model.GctfDataManage.Get(&u)
	if err != nil || h == false {
		log.Println("user/startProblem: error to search db(username)" + sp.Username + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to search db"})
	}
	h, err = model.GctfDataManage.Get(&p)
	if err != nil || h == false {
		log.Println("user/startProblem: error to search db(problem)" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to search db"})
	}
	up.UserId = u.Id
	up.ProblemsId = sp.Problem_ID
	//c,err:=gctfConfig.DockerClient.ContainerCreate(context.Background(),nil,nil,nil,"")
}
func startContainer(name string) error {
	createOpt := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: name,
		},
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
			PortBindings: map[docker.Port][]docker.PortBinding{
				"2817": {
					{
						"0.0.0.0",
						"",
					},
				},
			},
		},
	}
	rsp, err := gctfConfig.DockerClient.CreateContainer(createOpt)
}
func GetProblemList(c *gin.Context) {

}
