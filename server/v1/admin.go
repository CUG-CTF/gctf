package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"../gctfUtils"
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
func ChangeGCTFMode(c *gin.Context) {

}

func UploadProblem(c *gin.Context) {
	category:=c.PostForm("category")
	description:=c.PostForm("description")
	problem,err:=c.FormFile("problem")
	if err!=nil{
		log.Println("UploadProblem error:" +err.Error())
		c.JSON(http.StatusBadRequest,gin.H{"msg":"upload file erroe"})
		return
	}
	//Todo: write db, security check to vars
	category=category
	description=description

	err=os.MkdirAll("problems/"+category,os.ModePerm)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"msg":"error to create folder"+"problems/"+category,"err":err.Error()})
	}
	filename:=problem.Filename
	//TODO: add zip support
	err=c.SaveUploadedFile(problem,"problems/"+category+"/"+filename)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"msg":"error to save upload file","err":err.Error()})
	}
	f,err:=os.Open("problems/"+category+"/"+filename)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"msg":"error to read upload file","err":err.Error()})
	}
	gctfUtils.ExtractTarGz(f)

}
//change a problem's information
func ChangeProblem(c *gin.Context) {

}

//delete a problem from disk and db
func DeleteProblem(c *gin.Context) {

}