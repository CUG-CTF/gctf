package v1

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	. "github.com/CUG-CTF/gctf/server/model"
	"github.com/gin-gonic/gin"
	b "golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UserInfo(c *gin.Context) {
	type retInfo struct {
		Username       string    `json:"username" xorm:"unique pk"`
		Email          string    `json:"email" xorm:"unique"`
		RegisterTime   time.Time `json:"register_time"`
		SolvedProblems string    `json:"solved_problem"`
		Score          int       `json:"score"`
	}
	var u User
	var ui retInfo
	err := c.BindJSON(&u)
	//todo : error handle
	//todo:必须限定username，不然会查到别的数据
	h, err := GctfDataManage.Where("username =?", u.Username).Get(&u)
	if err != nil {
		log.Println("User/UserInfo :error to query db(username) ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error username"})
		return
	}
	if !h {
		log.Println("User/UserInfo: attempt to attach? no this user", u)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error Data"})
		return
	}
	ui.Username = u.Username
	ui.Email = u.Email
	ui.RegisterTime = u.RegisterTime
	ui.Score = u.Score
	ui.SolvedProblems = u.SolvedProblems
	c.JSON(http.StatusOK, &ui)

}

var Sessions map[string][]string

func WriteSession(username, token string) {
	//TODO: 应该用redis存session，这里存的session不会过期!!
	if Sessions == nil {
		Sessions = make(map[string][]string)
	}
	Sessions[username] = append(Sessions[username], token)
}

func checkSessionMiddleware(c *gin.Context) {
	//TODO: redis
	t := struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}{}
	//todo:gin 框架不能读两次body
	data, err := c.GetRawData()
	if err != nil {
		log.Println("checkSessionMiddleware: error to read data!", err.Error())
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Println("checkSessionMiddleware:error to bind json! ", err.Error())
	}
	if len(t.Username) == 0 || len(t.Token) == 0 {
		log.Printf("checkSessionMiddleware: wrong request! %v", t)
	}

	val, ok := Sessions[t.Username]
	if !ok {
		c.Redirect(http.StatusMovedPermanently, "/user/login")
		c.Abort()
	}
	//TODO:string compare is slow
	for _, x := range val {
		if t.Token == x {
			c.Next()
			return
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/user/login")
	c.Abort()
}

//
//TODO：检查是否为admin(查数据库)
func checkAdmin(c *gin.Context) {
	t := struct {
		Token    string `json:"token"`
		Usernmae string `json:"username"`
	}{}
	data, err := c.GetRawData()
	if err != nil {
		log.Println("checkSessionMiddleware: error to read data!", err.Error())
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	err = json.Unmarshal(data, &t)
	if t.Usernmae != "gctf" {
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
		Username string `json:"username"`
		Token    string `json:"token"`
		Message  string `json:"msg"`
	}
	var lr loginReturn
	var u User
	u.Username = l.User
	h, err := GctfDataManage.Where("username =?", u.Username).Get(&u)
	if err != nil {
		log.Println("user/login: error to get user info", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "login error"})
		return
	}
	if !h {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "login error"})
		return
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
		lr.Username = u.Username
		userToken, err := b.GenerateFromPassword([]byte(u.Username+time.Now().String()), b.DefaultCost)
		if err != nil {
			log.Println("user/login: error to gen token:", u, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to gen token"})
			return
		}
		lr.Token = base64.StdEncoding.EncodeToString(userToken)
		WriteSession(l.User, lr.Token)
		c.JSON(http.StatusOK, &lr)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"msg": "login error"})
		return
	}
}
func Logout(c *gin.Context) {
	//TODO: del session from K-V
	t := struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}{}
	//todo: error handle
	_ = c.BindJSON(&t)
	_, ok := Sessions[t.Username]
	if ok {
		delete(Sessions, t.Username)
	}
	c.JSON(http.StatusOK, gin.H{"msg": "logout ok"})

}

func Register(c *gin.Context) {
	var newUser User
	// username,password,email need
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "bad request" + err.Error()})
		return
	}
	//TODO:验证邮箱，用户名，密码格式
	hashed, err := b.GenerateFromPassword([]byte(newUser.Password), b.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error encrypt password!" + err.Error()})
		return
	}
	newUser.Password = base64.StdEncoding.EncodeToString(hashed)
	_, err = GctfDataManage.Insert(&newUser)
	if err != nil {
		log.Println("register error:" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "user existed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "OK"})
}

//TODO:获取分数
func GetScore(c *gin.Context) {

}

func SubmitFlag(c *gin.Context) {
	type submitFlag struct {
		Username   string `json:"username"`
		Problem_id string `json:"problem_id"`
		Flag       string `json:"flag"`
	}
	var myflag submitFlag
	err := c.BindJSON(&myflag)
	if err != nil {
		log.Println("user/SubmitFLag error to bind json", myflag, err)
		c.JSON(http.StatusBadRequest, gin.H{"succeed": false, "msg": "error to submit you flag"})
		return
	}
	var u User
	u.Username = myflag.Username
	h, err := GctfDataManage.Where("username =?", u.Username).Get(&u)
	//前端试图去提交一个不存在的用户名，并绕过了token!
	if !h {
		log.Println("user/SubmitFLag: attempt to attack? no this user!", u)
		c.JSON(http.StatusBadRequest, gin.H{"succeed": false, "msg": "error to submit you flag"})
		return
	}
	if err != nil {
		log.Println("user/SubmitFlag: error to query db(user)", err)
		c.JSON(http.StatusBadRequest, gin.H{"succeed": false, "msg": "error to submit you flag"})
		return
	}
	//solvedProblems := strings.Split(u.SolvedProblems, ",")
	//for _, solvedProblem := range solvedProblems {
	//	if myflag.Problem_id == solvedProblem {
	//		c.JSON(http.StatusOK, gin.H{"msg": "You submit already"})
	//		return
	//	}
	//}
	var p Problems
	p.Id, _ = strconv.ParseInt(myflag.Problem_id, 10, 64)
	//查database去拿到正确的flag
	h, err = GctfDataManage.ID(p.Id).Get(&p)
	if !h {
		log.Println("user/SubmitFLag: attempt to attack? no this problem_ID!", u)
		c.JSON(http.StatusBadRequest, gin.H{"succeed": false, "msg": "error to submit you flag"})
		return
	}
	if err != nil {
		log.Println("user/SubmitFlag: error to query db(problem)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"succeed": false, "msg": "error to submit you flag"})
		return
	}

	if GCTFConfig.GCTF_MODE {
		//Todo:在比赛模式中，应该重新计算所有人的分数

	} else {
		if p.Flag == myflag.Flag {
			if len(u.SolvedProblems) == 0 {
				//第一次提交flag，不然就逗号开头了
				u.SolvedProblems = myflag.Problem_id
			} else {
				solved := strings.Split(u.SolvedProblems, ",")
				for _, x := range solved {
					if x == myflag.Problem_id {
						c.JSON(http.StatusBadRequest, gin.H{"msg": "This flag already submit!"})
						return
					}
				}
				u.SolvedProblems += "," + myflag.Problem_id
			}

			//更新分数
			u.Score += p.Value
			n, err := GctfDataManage.Where("username =?", u.Username).Cols("score", "solved_problems").Update(&u)
			if n != 1 || err != nil {
				log.Println("user/SubmitFlag: error to update user db or update user data not only one! ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "server internal error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"succeed": true, "msg": "flag correct!"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"succeed": false, "msg": "flag error!"})
			return
		}
	}

}
