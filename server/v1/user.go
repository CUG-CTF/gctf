package v1

import (
	. "../model"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	b "golang.org/x/crypto/bcrypt"
)

func UserInfo(c *gin.Context) {
}
func Logout(c *gin.Context) {

}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	hashed, err := b.GenerateFromPassword([]byte(password), b.DefaultCost)
	if err != nil {
		c.Error(err)
		return
	}
	newUser := User{
		Username: username,
		Password: base64.StdEncoding.EncodeToString(hashed)
	}
}

func GetScore(c *gin.Context) {

}

func SubmitFlag(c *gin.Context) {

}

func UserDelProblem(c *gin.Context) {

}
