package v1

import (
	conf "../gctfConfig"
	"../model"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	b "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func UserInfo(c *gin.Context) {

}

var Sessions map[string][]string

func WriteSession(username, token string) {
	//TODO: check redis is available
	if Sessions == nil {
		Sessions = make(map[string][]string)
	}
	Sessions[username] = append(Sessions[username], token)
}

func checkSessionMiddleware(c *gin.Context) {
	//TODO: redis
	username, err := c.Cookie("username")
	if err != nil {
		log.Println("check session username:" + err.Error())
	}
	token, err := c.Cookie("token")
	if err != nil {
		log.Println("check session token:" + err.Error())
	}
	val, ok := Sessions[username]
	if !ok {
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
	}
	//string is slow
	for _, x := range val {
		if token == x {
			c.Next()
			return
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/login")
	c.Abort()
}
func Login(c *gin.Context) {
	//TODO :add json require
	type login struct {
		User     string `json:"username"`
		Password string `json:"password"`
	}
	var l login
	err := c.BindJSON(&l)
	if conf.GCTF_DEBUG {
		log.Println("user login:" + fmt.Sprintf("%#v", l))
	}
	if err != nil {
		log.Println("user login:" + err.Error())
	}
	//TODO:login check db
	type loginReturn struct {
		Message string `json:"message"`
		Token   string `json:"Token"`
	}
	var lr loginReturn

	if l.User == "gctf" && l.Password == "gctf" {
		lr.Message = "login ok"
		lr.Token = "gctf"
		WriteSession(l.User, lr.Token)
		c.JSON(http.StatusOK, &lr)
		c.SetCookie("username", l.User, 36000, "/", conf.GCTF_DOMAIN, false, true)
	}else {
		c.JSON(http.StatusForbidden,gin.H{"msg":"password error!"})
	}

}
func Logout(c *gin.Context) {
	//TODO: del session from K-V
}

func Register(c *gin.Context) {
	//TODO: add email verify,write in db
	var newUser model.User
	err:=c.BindJSON(&newUser)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"msg":"bad request"+err.Error()})
	}
	//username := c.PostForm("username")
	//password := c.PostForm("password")
	hashed, err := b.GenerateFromPassword([]byte(newUser.Password), b.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error encrypt password!"+err.Error()})
		return
	}
	newUser.Password=base64.StdEncoding.EncodeToString(hashed)
	_,err=model.GctfDataManage.Insert(&newUser)
	if err!=nil{
		log.Println("register error:"+err.Error())
		c.JSON(http.StatusInternalServerError,gin.H{"msg":"error to insert to db!"})
	}
	c.JSON(http.StatusOK,gin.H{"msg":"OK"})
}

func GetScore(c *gin.Context) {

}

func SubmitFlag(c *gin.Context) {

}

func UserDelProblem(c *gin.Context) {

}
