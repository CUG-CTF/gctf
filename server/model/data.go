package model

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"time"
)

type User struct {
	Id int64 `xorm:"autoincr 'id'"`
	//seem unique not work
	Username       string `xorm:"unique pk"`
	Password       string
	Email          string    `xorm:"unique"`
	RegisterTime   time.Time `xorm:"created notnull"`
	ProblemsId     string
	SolvedProblems string
	Score          int
	IsAdmin        bool `xorm:"'is_admin'"`
}

// Problem table, should Location is fixed?
//TODO:unique invalid?
type Problems struct {
	Id          int64  `xorm:"autoincr pk 'id'"`
	Name        string `xorm:"unique(name_location)"`
	Description string // Problem Description
	Value       int    // score
	Category    string
	Hidden      bool                                  // should be problem hide?
	Location    string `xorm:"unique(name_location)"` // saved physical position
	Flag        string                                //默认flag，(不开启动态flag)
	Scale       int `xorm:"default 0"`                // score scale when each answer submit
}

// 每启动一个problem 实例，就写入一条数据，如果flag为动态，那么就要填入Flag字段
type UserProblems struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	UserId     int64 `xorm:"unique(user_problem)"` //foreignkey gctf_user.id
	Location   string                              // problem net location
	Flag       string
	ProblemsId int64 `xorm:"unique(user_problem)"` //foreignkey gctf_problems.id
	Expired    time.Time                           //过期时间
}

type Hints struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	ProblemsId int64 // foreignkey gctf_problems.id
	Hint       string
	Cost       int // cost score to get hint
}

//type Tag struct {
//	Id         int64 `xorm:"autoincr pk 'id'"`
//	ProblemsId int64 //foreignkey gctf_problems.id
//	Tag        string
//}

type Teams struct {
	Id     int64  `xorm:"autoincr pk 'id'"`
	Name   string `xorm:"unique"`
	Member string
	Banned bool   // if true this team can't login
	Token  string // team token
}

type GCTFConfigStruct struct {
	GCTF_PORT            string `json:"port"`
	GCTF_DEBUG           bool   `json:"debug"`
	GCTF_MODE            bool   `json:"mode"` //true is contest
	GCTF_PROBLEM_TIMEOUT int    `json:"problem_create_timeout"`
	GCTF_EXPLIRED_TIME   int    `json:"expired_time"` // 过期时间，单位分钟
	GCTF_DB_DRIVER       string `json:"database_type"`
	GCTF_DB_STRING       string `json:"database_address"`
	GCTF_DOMAIN          string `json:"domain_name"`
	//TODO: add docker server manager,else use local docker unix sock
	GCTF_DOCKERS []string `json:"docker_servers"`
}
type DockerManager interface {
	GetDockerClient() *docker.Client
}

var (
	GCTFConfig *GCTFConfigStruct
	//TODO: add docker server manager,else use local docker unix sock
	//only in dev
	GCTFDockerManager DockerManager
)

//database manager
var GctfDataManage *xorm.Engine
