package v1

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/CUG-CTF/gctf/server/model"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	b "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//TODO:添加用户，直接调用Register接口，缺少测试
func AddUsers(c *gin.Context) {
	//username password
	Register(c)
}

func ChangeUserPasswd(c *gin.Context) {
	user := struct {
		Username string
		Password string
	}{}
	err := c.BindJSON(&user)
	if err != nil {
		log.Println("admin/ChangeUser")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to bind json!"})
		return
	}
	var u model.User
	u.Username = user.Username
	_, err = model.GctfDataManage.Get(&u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to search db"})
		return
	}
	hashed, err := b.GenerateFromPassword([]byte(user.Password), b.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error encrypt password!" + err.Error()})
		return
	}
	u.Password = base64.StdEncoding.EncodeToString(hashed)
	n, err := model.GctfDataManage.Id(u.Id).Cols("Password").Update(&u)
	if n != 1 {
		log.Println("admin/changeUserpasswd change passwd error,not only one!:" + fmt.Sprintf("%v", user))
	}
	c.JSON(http.StatusOK, gin.H{"msg": "change user passwd success"})
	return

}

// 单独设置一个user的题目
func SetUserProblem(c *gin.Context) {

}

//随机出题给所有用户，应该的比赛模式中使用
func RandomAllUsersProblem(c *gin.Context) {

}

// 比赛或者训练模式
func ChangeGCTFMode(c *gin.Context) {

}

//admin upload a problem which include Dockerfile,  only support tar,tar.gz format
func UploadProblem(c *gin.Context) {
	//得到form post的参数

	name := c.PostForm("name")
	category := c.PostForm("category")
	description := c.PostForm("description")
	form_value := c.PostForm("value")
	flag := c.PostForm("flag")
	form_port := c.PostForm("port")
	//传过来的tar 或者tar.gz
	problem, err := c.FormFile("problem")
	if err != nil {
		log.Println("Upload problem error:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "upload file error"})
		return
	}
	//题目分数应该是整数
	value, err := strconv.Atoi(form_value)
	if err != nil {
		log.Println("Upload problem's value error:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "wrong  value!"})
		return
	}
	if value < 0 {
		log.Println("Upload problem's value not a negative number")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "negative value"})
		return
	}
	problemPort, err := strconv.Atoi(form_port)
	if err != nil {
		log.Println("Upload problem's port's value error: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "wrong port!"})
		return
	}
	if problemPort >= 65535 {
		log.Println("Upload problem's port's value too large!")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "too large port"})
		return
	}

	var p model.Problems
	p.Name = name
	h, err := model.GctfDataManage.Get(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to search db! " + err.Error()})
		return
	}
	if h {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "name exist!"})
		return
	}

	//Todo: 对post过来的参数进行安全检查和过滤
	err = os.MkdirAll("problems/"+category, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to create folder" + "problems/" + category, "err": err.Error()})
	}
	dist:="problems/"+category+"/"+name+".tar.gz"
	_,err=os.Stat(dist)
	if err==nil{
		c.JSON(http.StatusOK,gin.H{"msg":"file exist!"})
		return
	}
	err = c.SaveUploadedFile(problem, dist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to save upload file", "err": err.Error()})
		return
	}
	p.Category = category
	p.Value = value
	p.Description = description
	p.Name = name
	p.Hidden = false
	p.Flag = flag
	p.Location = "problems/" + category + "/" + name + ".tar.gz"
	p.Port = problemPort
	// 加入队列
	go func() { model.UploadQuene.UploadProblems <- p }()
	//启动上传协程
	if !model.UploadQuene.BuildBegin {
		go startBuildProblem()
		model.UploadQuene.BuildBegin = true
	}

	c.JSON(http.StatusOK, gin.H{"msg": "start building"})
	return

}

func startBuildProblem() {
	for p := range model.UploadQuene.UploadProblems {
		go buildUploadProblem(p)
	}
}

func buildUploadProblem(p model.Problems) {
	//TODO:!似乎有问题
	//TODO: context必须设置一个timeout
	//TODO:在所有docker server上创建题目,构建后进行pull，两种方案:
	// 1. master作为docker hub，pull from master
	// 2. master docker hub 分离，pull from docker hub
	var br model.BuildResult
	f, err := os.Open("problems/" + p.Category + "/" + p.Name)
	if err!=nil{
		log.Println("admin/Upload_Problem:error to open problem!",err.Error())
		return
	}
	//在数据库，文件层面已经做了一次problem name 重复check了，这里应该没必要做了吧- -
	buildTimeLimit := model.GCTFConfig.GCTF_BUILD_TIME_LIMIT
	dockerContext, _ := context.WithTimeout(context.Background(), time.Duration(buildTimeLimit)*time.Minute)
	br.Output = bytes.NewBuffer(nil)
	bo := docker.BuildImageOptions{
		Context:      dockerContext,
		InputStream:  f,
		OutputStream: br.Output,
		Name:         p.Name,
	}
	cli := model.GCTFDockerManager.GetDockerClient()
	err = cli.BuildImage(bo)
	model.BuildOutputList.Add(p.Name, &br)
	if err != nil {
		br.Err = err
		log.Println("admin/Upload Problem:error to build a problem",err.Error())
		return
	}
	//存在已经判定过了
	_, err = model.GctfDataManage.Insert(&p)
	if err != nil {
		log.Println("db problem insert error:", err.Error())
		return
	}

}

//TODO:change a problem's information
//验证数据库、镜像、磁盘
//value flag、description、category、Hidden、Port
func ChangeProblem(c *gin.Context) {

}

//TODO:delete a problem from disk and db
//先查db，然后删镜像，最后删磁盘文件
func DeleteProblem(c *gin.Context) {
	var pn model.Problems
	err:=c.BindJSON(&pn)
	if err!=nil{
		c.JSON(http.StatusBadRequest,"error to bind json!")
		log.Println("admin/DeleteProblem: error to bind json: "+err.Error())
		return
	}
	h,err:=model.GctfDataManage.Get(&pn)
	_=h
	if err!=nil{
		c.JSON(http.StatusInternalServerError,"errro to search db!")
		log.Println("admin/DeleteProblem: error to search db: "+err.Error())
		return
	}
	cli:=model.GCTFDockerManager.GetDockerClient()
	cli.RemoveContainer(docker.RemoveContainerOptions{})
}

//TODO:增加备份和恢复功能
//TODO:增减校验：数据库，硬盘，镜像三者同步
