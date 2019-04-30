package model

import (
	"bytes"
	"github.com/fsouza/go-dockerclient"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"sync"
	"time"
)

type User struct {
	Id int64 `json:"id" xorm:"autoincr 'id'"`
	//seem unique not work
	Username       string    `json:"username" xorm:"unique pk"`
	Password       string    `json:"password"`
	Email          string    `json:"email" xorm:"unique"`
	RegisterTime   time.Time `json:"register_time" xorm:"created notnull"`
	ProblemsId     string
	SolvedProblems string `json:"SolvedProblem"`
	Score          int    `json:"score"`
	IsAdmin        bool   `xorm:"'is_admin'"`
}

// Problem table, should Location is fixed?
//TODO:unique invalid?
type Problems struct {
	Id          int64  `xorm:"autoincr pk 'id'"`
	Name        string `xorm:"unique"`
	Description string // Problem Description
	Value       int    // score
	Category    string
	Hidden      bool                   // should be problem hide?
	Location    string `xorm:"unique"` // saved physical position
	Flag        string                 //默认flag，(不开启动态flag)
	Scale       int `xorm:"default 0"` // score scale when each answer submit
	Port        int `xorm:"default 0"` //内部端口
}

// 每启动一个problem 实例，就写入一条数据，如果flag为动态，那么就要填入Flag字段
type UserProblems struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	UserId     int64 `xorm:"unique(user_problem)"` //foreignkey gctf_user.id
	Location   string                              // problem net location
	Flag       string
	ProblemsId int64 `xorm:"unique(user_problem)"` //foreignkey gctf_problems.id
	DockerID   string                              //docker id
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
	GCTF_PORT             string   `json:"port"`
	GCTF_DEBUG            bool     `json:"debug"`
	GCTF_MODE             bool     `json:"mode"` //true is contest
	GCTF_PROBLEM_TIMEOUT  int      `json:"problem_create_timeout"`
	GCTF_EXPLIRED_TIME    int      `json:"expired_time"` // 过期时间，单位分钟
	GCTF_BUILD_TIME_LIMIT int      `json:"build_time_limit"`
	GCTF_DB_DRIVER        string   `json:"database_type"`
	GCTF_DB_STRING        string   `json:"database_address"`
	GCTF_DOMAIN           string   `json:"domain_name"`
	GCTF_DOCKERS          map[string]string `json:"docker_servers"`
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

//problems 上传队列
type uploadQueue struct {
	m              sync.Mutex
	UploadProblems chan Problems
	BuildBegin     bool
}

//problems 构建输出列表，不要直接操作BuildOutputs
type buildOutputList struct {
	m            sync.Mutex
	buildOutputs map[string]*BuildResult
}

//problems  构建结果
type BuildResult struct {
	Err    error
	Output *bytes.Buffer
}

var (
	UploadQuene     uploadQueue
	BuildOutputList buildOutputList
)

//添加problem build 输出，并不会检查是否存在，需要用Get自行判断
//TODO: 增加定时清空
func (b buildOutputList) Add(name string, buildResult *BuildResult) {
	b.m.Lock()
	b.buildOutputs[name] = buildResult
	b.m.Unlock()
}

//获取一个problem build输出
func (b buildOutputList) Get(name string) *BuildResult {
	b.m.Lock()
	r, ok := b.buildOutputs[name]
	b.m.Unlock()
	if ok {
		return r
	}
	return nil
}

func initBuildQueue() {
	UploadQuene.BuildBegin = false
	UploadQuene.UploadProblems = make(chan Problems, 10)
}
func init() {
	BuildOutputList.buildOutputs = make(map[string]*BuildResult)
	initBuildQueue()
}
