package v1

import (
	. "../model"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	b "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UserInfo(c *gin.Context) {

}

var Sessions map[string][]string

func WriteSession(username, token string) {
	//TODO: 应该用redis存session，这里存的session不会过期
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

//
//TODO：检查是否为admin(查数据库)
func checkAdmin(c *gin.Context) {
	username, _ := c.Cookie("username")
	if username != "gctf" {
		c.JSON(http.StatusForbidden, gin.H{"msg": "you are not admin"})
		c.Abort()
	}
	c.Next()
}
func Login(c *gin.Context) {
	//TODO :增加json校验，不允许为空
	type login struct {
		User     string `json:"username"`
		Password string `json:"password"`
	}
	var l login
	err := c.BindJSON(&l)
	if GCTFConfig.GCTF_DEBUG {
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
	var u User
	u.Username = l.User
	h, err := GctfDataManage.Get(&u)
	if err != nil {
		log.Println("user/login: error to get user info", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "login error"})
	}
	if !h {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "login error"})
	}
	hashPasswd, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		log.Println("user/login:error to decode passwd", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "login error"})
		return
	}

	err = b.CompareHashAndPassword(hashPasswd, []byte(l.Password))
	if err == nil {
		lr.Message = "login ok"

		userToken, err := b.GenerateFromPassword([]byte(u.Username+time.Now().String()), b.DefaultCost)
		if err != nil {
			log.Println("user/login: error to gen token:", u, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to gen token"})
			return
		}
		lr.Token = base64.StdEncoding.EncodeToString(userToken)
		WriteSession(l.User, lr.Token)
		c.SetCookie("username", l.User, 36000, "/", GCTFConfig.GCTF_DOMAIN, false, true)
		c.SetCookie("token", lr.Token, 36000, "/", GCTFConfig.GCTF_DOMAIN, false, true)
		c.JSON(http.StatusOK, &lr)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"msg": "login error"})
	}
}
func Logout(c *gin.Context) {
	//TODO: del session from K-V
}

func Register(c *gin.Context) {
	//TODO: add email verify
	var newUser User
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad request" + err.Error()})
		return
	}
	//username := c.PostForm("username")
	//password := c.PostForm("password")
	hashed, err := b.GenerateFromPassword([]byte(newUser.Password), b.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error encrypt password!" + err.Error()})
		return
	}
	newUser.Password = base64.StdEncoding.EncodeToString(hashed)
	_, err = GctfDataManage.Insert(&newUser)
	if err != nil {
		log.Println("register error:" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to insert to db!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "OK"})
}

//TODO:获取分数
func GetScore(c *gin.Context) {

}

func SubmitFlag(c *gin.Context) {
	type submitFlag struct {
		username   string
		problem_id string
		flag       string
	}
	var myflag submitFlag
	err := c.BindJSON(&myflag)
	if err != nil {
		log.Println("user/SubmitFLag error to bind json", myflag, err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to submit you flag"})
		return
	}
	var u User
	u.Username = myflag.username
	h, err := GctfDataManage.Get(&u)
	//前端试图去提交一个不存在的用户名，并绕过了token!
	if !h {
		log.Println("user/SubmitFLag: attempt to attack? no this user!", u)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to submit you flag"})
		return
	}
	if err != nil {
		log.Println("user/SubmitFlag: error to query db(user)", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to submit you flag"})
		return
	}
	solvedProblems := strings.Split(u.SolvedProblems, ",")
	for _, solvedProblem := range solvedProblems {
		if myflag.problem_id == solvedProblem {
			c.JSON(http.StatusOK, gin.H{"msg": "You submit already"})
			return
		}
	}
	var p Problems
	p.Id, _ = strconv.ParseInt(myflag.problem_id, 10, 64)
	//查database去拿到正确的flag
	h, err = GctfDataManage.Get(&p)
	if !h {
		log.Println("user/SubmitFLag: attempt to attack? no this problem_ID!", u)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error to submit you flag"})
		return
	}
	if err != nil {
		log.Println("user/SubmitFlag: error to query db(problem)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to submit you flag"})
		return
	}

	if GCTFConfig.GCTF_MODE {
		//Todo:在比赛模式中，应该重新计算分数

	} else {
		if p.Flag == myflag.flag {
			if len(u.SolvedProblems) == 0 {
				//第一次提交flag，不然就逗号开头了
				u.SolvedProblems = myflag.problem_id
			}
			u.SolvedProblems += "," + myflag.problem_id

			//更新分数
			u.Score += p.Value
			n, err := GctfDataManage.Cols("score", "solved_problems").Update(&u)
			if n != 1 || err != nil {
				log.Println("user/SubmitFlag: error to update user db or update user data not only one! ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "server internal error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"msg": "flag correct!"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": "flag error!"})
			return
		}
	}

}

//TODO:用户删除题目实例
func UserDelProblem(c *gin.Context) {

}
