package v1

import (
	"../model"
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func StartProblem(c *gin.Context) {
	//TODO: 增加更多的debug输出信息
	expired_time := model.GCTFConfig.GCTF_EXPLIRED_TIME
	//计算过期时间
	now := time.Now().Add(time.Duration(expired_time) * time.Minute)
	// POST DATA
	type UserStartProblem struct {
		Username   string `json:"username"`
		Token      string `json:"token"`
		Problem_id int64  `json:"problem_id"`
	}
	var sp UserStartProblem
	// json 反序列化失败
	err := c.BindJSON(&sp)
	if err != nil {
		log.Println("user/StartProblem: errot to bind json" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to start bind json"})
		return
	}
	var p model.Problems
	var u model.User
	u.Username = sp.Username
	// 得到user ID
	h, err := model.GctfDataManage.Get(&u)
	if err != nil || h == false {
		log.Println("user/StartProblem: error to search db(username)" + sp.Username + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to search db"})
		return
	}

	p.Id = sp.Problem_id
	//得到problem_id
	h, err = model.GctfDataManage.Get(&p)
	if err != nil || h == false {
		log.Println("user/StartProblem: error to search db(problem) " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to search db"})
		return
	}

	var up model.UserProblems
	up.ProblemsId = sp.Problem_id
	up.UserId = u.Id
	//查一下是不是已经创建题目实例了

	h, err = model.GctfDataManage.Get(&up)
	if err != nil {
		log.Println("user/StartProblem: error to search db(user problem) ", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to search db"})
		return
	}
	if h {
		//Todo: 检查过期时间，可能已经过期但未删除
		host_port_ip := strings.Split(up.Location, ":")
		if len(host_port_ip) != 2 {
			//Todo: 应该还有别的处理（删掉容器和这条数据）
			log.Println("user/StartProblem: error Location format!")
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to start problem"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"host_ip":   host_port_ip[0],
			"host_port": host_port_ip[1],
			"expired:":  up.Expired.Format("15:04:05"),
		})
		return
	}
	//启动实例
	problemAddr, err := startContainer(p.Name)
	//TODO:启动失败，应当删除题目实例
	if err != nil {
		log.Println("user/StartProblem: error to start a  problem(name =" + p.Name + ") " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to start a problem"})
		return
	}
	//返回题目地址
	up.Location = problemAddr.HostIP + ":" + problemAddr.HostPort
	up.UserId = u.Id
	up.ProblemsId = p.Id
	_, err = model.GctfDataManage.Insert(&up)
	if err != nil {
		log.Printf("user/StartProblem error to insert to db: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to insert db"})
	}
	c.JSON(http.StatusOK, gin.H{
		"host_ip":   problemAddr.HostIP,
		"host_port": problemAddr.HostPort,
		"expired:":  now.Format("15:04:05"),
	})
}
func startContainer(name string) (*docker.PortBinding, error) {
	//TODO:设置1分钟测试用，实际开发要替换为配置文件中设置的时间
	context_timeout, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	//context_timeout,_:=context.WithTimeout(context.Background(),time.Duration(model.GCTFConfig.GCTF_PROBLEM_TIMEOUT)*time.Minute)
	createOpt := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: name,
		},
		Context: context_timeout,
		//TODO:web题目多端口处理
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
	if model.GCTFConfig.GCTF_MODE {
		log.Println("attempt to read all problem!")
		c.JSON(http.StatusForbidden, gin.H{"msg": "not allow"})
		return
	}
	type retData struct {
		Id          int64  `json:"id"`
		Name        string `json:"nane"`
		Description string `json:"decription"`
		Value       int    `json:"value"`
		Category    string `json:"category"`
	}
	var problems []model.Problems
	var retList []retData
	err := model.GctfDataManage.Find(&problems)
	if err != nil {
		log.Println("problem/GetProblemList:error to get all problems", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to get problem list"})
		return
	}
	for _, x := range problems {
		if !x.Hidden {
			retList = append(retList, retData{
				x.Id,
				x.Name,
				x.Description,
				x.Value,
				x.Category,
			})
		}
	}
	c.JSON(http.StatusOK, retList)
}
