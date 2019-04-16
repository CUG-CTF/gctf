package v1

import (
	"github.com/CUG-CTF/gctf/server/model"
	"bytes"
	"context"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//TODO:添加用户
func AddUsers(c *gin.Context) {

}

//TODO:改密码
func ChangeUserPasswd(c *gin.Context) {

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
		c.JSON(http.StatusBadRequest, gin.H{"msg": "wrong value!" + err.Error()})
	}

	//Todo: 对post过来的参数进行安全检查和过滤
	err = os.MkdirAll("problems/"+category, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to create folder" + "problems/" + category, "err": err.Error()})
	}
	filename := problem.Filename
	err = c.SaveUploadedFile(problem, "problems/"+category+"/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to save upload file", "err": err.Error()})
	}
	f, err := os.Open("problems/" + category + "/" + filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to read upload file", "err": err.Error()})
	}
	build_result := buildUploadProblem(f, name)
	if build_result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to build image: " + build_result.Error()})
		return
	}
	var p model.Problems
	p.Category = category
	p.Value = value
	p.Description = description
	p.Name = name
	p.Hidden = false
	p.Flag = flag
	p.Location = "problems/" + category + "/" + filename
	_, err = model.GctfDataManage.Insert(p)
	if err != nil {
		log.Println("db problem insert error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error when insert to db!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "upload succeed!"})
}

//TODO:change a problem's information
func ChangeProblem(c *gin.Context) {

}

//TODO:delete a problem from disk and db
//先查db，然后删镜像，最后删磁盘文件
func DeleteProblem(c *gin.Context) {

	model.GctfDataManage.Query()
}

func buildUploadProblem(f io.Reader, name string) error {
	//TODO: context必须设置一个timeout
	//TODO:在所有docker server上创建题目
	buildTimeLimit:=model.GCTFConfig.GCTF_BUILD_TIME_LIMIT
	dockerContext,_ := context.WithTimeout(context.Background(),time.Duration(buildTimeLimit)*time.Minute)
	buildOutput := bytes.NewBuffer(nil)
	bo := docker.BuildImageOptions{
		Context:      dockerContext,
		InputStream:  f,
		OutputStream: buildOutput,
		Name:         name,
	}
	cli := model.GCTFDockerManager.GetDockerClient()
	err := cli.BuildImage(bo)
	if err != nil {
		log.Println("admin/Upload Problem:error to build a problem")
		return err
	}
	_, err = ioutil.ReadAll(buildOutput)

	//for true {
	//	_, err = buildOutput.Read(nil)
	//	if err == io.EOF {
	//		return nil
	//	}
	//	if err != nil {
	//		return err
	//	}
	//
	//}

	return err
}

//TODO:增加备份和恢复功能