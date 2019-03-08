package v1

import (
	"../gctfConfig"
	"../model"
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
)

//TODO: problems sync with database
func AddUsers(c *gin.Context) {

}
func ChangeUserPasswd(c *gin.Context) {

}

func SetUserProblem(c *gin.Context) {

}
func RandomAllUsersProblem(c *gin.Context) {

}
// 比赛或者训练
func ChangeGCTFMode(c *gin.Context) {

}
//admin upload a problem which include Dockerfile,  only support tar,tar.gz format
func UploadProblem(c *gin.Context) {
	name := c.PostForm("name")
	category := c.PostForm("category")
	description := c.PostForm("description")
	form_value := c.PostForm("value")
	problem, err := c.FormFile("problem")
	if err != nil {
		log.Println("Upload problem error:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "upload file erroe"})
		return
	}
	value, err := strconv.Atoi(form_value)
	if err != nil {
		log.Println("Upload problem's value error:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg":"wrong value!"+err.Error()})
	}

	//Todo: security check to vars

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
	build_result := buildUploadProblem(f,name)
	if build_result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to build image: " + build_result.Error()})
	}
	var p model.Problems
	p.Category = category
	p.Value = value
	p.Description = description
	p.Name = name
	p.Hidden=false
	p.Location="problems/" + category + "/" + filename
	_, err = model.GctfDataManage.Insert(p)
	if err != nil {
		log.Println("db problem insert error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error when insert to db!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "upload succeed!"})
}

//change a problem's information
func ChangeProblem(c *gin.Context) {

}

//delete a problem from disk and db
func DeleteProblem(c *gin.Context) {

}

func buildUploadProblem(f io.Reader,name string) error {
	//TODO: This context must set timeout
	dockerContext := context.Background()
	buildOutput:=bytes.NewBuffer(nil)
	bo:=docker.BuildImageOptions{
		Context:dockerContext,
		InputStream:f,
		OutputStream:buildOutput,
		Name:name,
	}
	err := gctfConfig.DockerClient.BuildImage(bo)
	if err != nil {
		log.Println("admin/Upload Problem:error to build a problem")
		return err
	}
	_,err=ioutil.ReadAll(buildOutput)

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
