package v1

import (
	"../model"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func StartProblem(c *gin.Context) {
	//TODO: add more debug msg
	now:=time.Now().Add(10*time.Minute)
	type UserStartProblem struct {
		Username   string `json:"username"`
		Token      string `json:"token"`
		Problem_ID int64  `json:"problem_id"`
	}
	var sp UserStartProblem
	err := c.BindJSON(&sp)
	if err != nil {
		log.Println("user/StartProblem: errot to bind json" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to start bind json"})
	}
	var p model.Problems
	var u model.User
	u.Username = sp.Username
	h, err := model.GctfDataManage.Get(&u)
	if err != nil || h == false {
		log.Println("user/StartProblem: error to search db(username)" + sp.Username + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to search db"})
	}
	p.Id = sp.Problem_ID
	h, err = model.GctfDataManage.Get(&p)
	if err != nil || h == false {
		log.Println("user/StartProblem: error to search db(problem)" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to search db"})
	}
	problemAddr, err := startContainer(p.Name)
	if err != nil {
		log.Println("user/StartProblem: error to start a  problem(name =" + p.Name + ") " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to start a problem"})
	}
	var up model.UserProblems
	up.Location=problemAddr.HostIP+":"+problemAddr.HostPort
	up.UserId=u.Id
	up.ProblemsId=p.Id
	_,err=model.GctfDataManage.Insert(&up)
	if err!=nil{
		log.Printf("user/StartProblem error to insert to db: "+err.Error())
		c.JSON(http.StatusInternalServerError,gin.H{"msg":"error to insert db"})
	}
	c.JSON(http.StatusOK,gin.H{
		"host_ip":problemAddr.HostIP,
		"host_port":problemAddr.HostPort,
		"expired:":now.Format("15:04:05"),
	})
}
func startContainer(name string) (*docker.PortBinding, error) {
	createOpt := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: name,
		},
		//TODO:dynamic port(s)?
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
	cli := model.GCTFDockerManager.GetDockerClient()
	rsp, err := cli.CreateContainer(createOpt)
	if err != nil {
		return nil, err
	}
	err = cli.StartContainer(rsp.ID, nil)
	if err != nil {
		return nil, err
	}
	rsp, err = cli.InspectContainer(rsp.ID)
	if err != nil {
		return nil, err
	}
	ret := new(docker.PortBinding)
	ret.HostIP = cli.Endpoint()
	ret.HostPort = rsp.NetworkSettings.Ports["2817"][0].HostPort
	return ret, err

}
func GetProblemList(c *gin.Context) {

}
