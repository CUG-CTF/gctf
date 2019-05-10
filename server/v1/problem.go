package v1

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/CUG-CTF/gctf/server/model"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
		Token      string `json:"token"`
		Username   string `json:"username"`
		Problem_id int64  `json:"problem_id"`
	}
	var sp UserStartProblem
	// json 反序列化失败
	err := c.BindJSON(&sp)
	if err != nil {
		log.Println("user/StartProblem: error to bind json" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to  bind json"})
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
	up.ProblemId = sp.Problem_id
	up.UserId = u.Id

	//查一下是不是已经创建题目实例了
	h, err = model.GctfDataManage.Get(&up)
	//TODO:更多测试出错原因
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
	problemAddr, id, cli, err := startContainer(p)
	if err != nil || len(id) == 0 {
		log.Println("user/StartProblem: error to start a  problem(name =" + p.Name + ") " + err.Error())
		_ = cli.RemoveContainer(docker.RemoveContainerOptions{ID: id, Force: true})
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to start a problem"})
		return
	}
	//生成动态flag

	//todo:fmt性能有点低,必须动态flag？
	randomBytes := []byte(fmt.Sprintf("%x", time.Now().UnixNano()))
	randomBytes = append(randomBytes, []byte(u.Username)...)
	randomFlag := fmt.Sprintf("%x", md5.Sum(randomBytes))
	//todo:可以更改的前缀
	randomFlag = "gctf{" + randomFlag + "}"
	changeFlagExec, err := cli.CreateExec(docker.CreateExecOptions{Container: id, Cmd: []string{"sh", "/changeFlag.sh", randomFlag}})
	if err == nil && changeFlagExec != nil {
		_, err = cli.StartExecNonBlocking(changeFlagExec.ID, docker.StartExecOptions{Detach: false, Tty: false})

	}
	if err != nil {
		log.Printf("user/StartProblem error to random flag!"+err.Error()+"%v /n", sp)
		randomFlag = ""
		err = nil
	}
	//返回题目地址
	up.Location = model.GCTFConfig.GCTF_DOCKERS[problemAddr.HostIP] + ":" + problemAddr.HostPort
	up.UserId = u.Id
	up.ProblemId = p.Id
	up.Expired = now
	up.DockerID = id
	up.Flag = randomFlag
	//Where("user_id=?",up.UserId).Where("problem_id=?",p.Id).
	_, err = model.GctfDataManage.Insert(&up)
	if err != nil {
		log.Printf("user/StartProblem error to insert to db: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to insert db"})
	}
	c.JSON(http.StatusOK, gin.H{
		"host_ip":   model.GCTFConfig.GCTF_DOCKERS[problemAddr.HostIP],
		"host_port": problemAddr.HostPort,
		"expired":   now.Format("15:04:05"),
	})
}
func startContainer(p model.Problems) (*docker.PortBinding, string, *docker.Client, error) {
	//TODO:设置10分钟测试用，实际开发要替换为配置文件中设置的时间
	context_timeout, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	//context_timeout,_:=context.WithTimeout(context.Background(),time.Duration(model.GCTFConfig.GCTF_PROBLEM_TIMEOUT)*time.Minute)
	createOpt := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: p.Name,
		},
		Context: context_timeout,
		//TODO:web题目多端口处理
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(strconv.Itoa(p.Port)): {
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
		return nil, "", cli, err
	}
	err = cli.StartContainer(rsp.ID, nil)
	if err != nil {
		return nil, "", cli, err
	}
	rsp, err = cli.InspectContainer(rsp.ID)
	if err != nil {
		return nil, "", cli, err
	}
	id := rsp.ID
	ret := new(docker.PortBinding)
	ret.HostIP = cli.Endpoint()
	//TODO:目前仅支持TCP端口
	port := docker.Port(strconv.Itoa(p.Port) + "/tcp")
	ret.HostPort = rsp.NetworkSettings.Ports[port][0].HostPort
	return ret, id, cli, err
}

func GetProblemList(c *gin.Context) {
	if model.GCTFConfig.GCTF_MODE {
		log.Println("attempt to read all problem!")
		c.JSON(http.StatusForbidden, gin.H{"msg": "not allow"})
		return
	}
	t := struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}{}
	err := c.BindJSON(&t)
	if err != nil {
		log.Println("user/GetProblemList error to bind json", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "bad request"})
		return
	}
	type problemList struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Value       int    `json:"value"`
		Category    string `json:"category"`
	}
	type retUserData struct {
		Categories  []string      `json:"categories"`
		ProblemList []problemList `json:"problem_list"`
	}
	type retAdminData struct {
		Categories  []string         `json:"categories"`
		ProblemList []model.Problems `json:"problem_list"`
	}
	var ru retUserData
	var ra retAdminData
	var problems []model.Problems
	var retList []problemList
	err = model.GctfDataManage.Find(&problems)
	if err != nil {
		log.Println("problem/GetProblemList:error to get all problems", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to get problem list"})
		return
	}
	ru.Categories = []string{"pwn", "web", "crypto", "re", "misc"}
	ra.Categories = []string{"pwn", "web", "crypto", "re", "misc"}
	ra.ProblemList = problems
	//管理员就获得所有题目
	if t.Username == "gctf" {
		c.JSON(http.StatusOK, ra)
		return
	}
	for _, x := range problems {
		if !x.Hidden {
			retList = append(retList, problemList{
				x.Id,
				x.Name,
				x.Description,
				x.Value,
				x.Category,
			})
		}
	}
	ru.ProblemList = retList
	c.JSON(http.StatusOK, ru)
}

//TODO:用户删除题目实例
//删掉容器，清除数据库
func UserDelProblem(c *gin.Context) {
	ud := struct {
		Token     string `json:"token"`
		Username  string `json:"username"`
		ProblemId int64  `json:"problem_id"`
	}{}
	err := c.BindJSON(&ud)
	if err != nil {
		log.Println("error to bind json:", err.Error())
		c.JSON(http.StatusBadRequest, "error to bind json!")
		return
	}
	var up model.UserProblems
	up.ProblemId = ud.ProblemId
	h, err := model.GctfDataManage.Get(&up)
	if err != nil {
		log.Println("problem/UserDelProblem : database error ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to search in db"})
		return
	}
	if !h {
		log.Println("problem/UserDelProblem attempt to del not exist problem! username: ", ud.Username)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "no this problem!"})
		return
	}
	_, err = model.GctfDataManage.Delete(&up)
	//todo:!client polling!
	cli := model.GCTFDockerManager.GetDockerClient()
	err = cli.RemoveContainer(docker.RemoveContainerOptions{ID: up.DockerID, Force: true})
	if err != nil {
		log.Println("problem/UserDelProblem error to del container! ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to del container!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "del container success!"})

}
