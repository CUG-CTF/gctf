package v1

import (
	"github.com/CUG-CTF/gctf/server/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strconv"
)

/*
[
{
"score": "15",
"username": "cug2"
},
{
"score": "10",
"username": "cug1"
}
]
 */
//User mode
func GetUsersRank(c *gin.Context) {
	data, err := model.GctfDataManage.Query("select `username`,`score` from gctf_user limit 100 ")
	var userName_scores []map[string]string
	for x := range data {
		userName_score := make(map[string]string)
		userName_score["username"] = string(data[x]["username"])
		userName_score["score"] = string(data[x]["score"])
		userName_scores = append(userName_scores, userName_score)
	}
	//d,_:=json.Marshal(userName_scores)
	if model.GCTFConfig.GCTF_DEBUG && err != nil {
		log.Println("public/GetUserRank:error to get userscore", data)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error to get user score"})
		return
	}
	//排序
	sort.Slice(userName_scores, func(i, j int) bool {
		a,_:=strconv.Atoi(userName_scores[i]["score"])
		b,_:=strconv.Atoi(userName_scores[j]["score"])
		return a>b
	})
	c.JSON(http.StatusOK, userName_scores)
}

//Team Mode?
func GetTeamsRank(c *gin.Context) {

}

//ping
func GctfPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "pong!"})
}
