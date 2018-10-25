package main

import (
	"../gctfConfig"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var GctfDataManage *xorm.Engine

func init() {
	var err error

	GctfDataManage, err = xorm.NewEngine(gctfConfig.GCTF_DB_NAME, gctfConfig.GCTF_DB_STRING)
	//All table name have gctf_ prefix
	tbMapper:=core.NewPrefixMapper(core.GonicMapper{}, "gctf_")
	// fix problem_I_D to problem_id
	GctfDataManage.SetColumnMapper(core.GonicMapper{})
	GctfDataManage.SetTableMapper(tbMapper)
	if err != nil {
		log.Fatal("database connet error:", err.Error())
	}
}

func main() {
	err:=GctfDataManage.Ping()
	if err!=nil{
		log.Fatal("error",err.Error())
	}
	if gctfConfig.GCTF_DEBUG{
		GctfDataManage.ShowSQL(true)
		GctfDataManage.Logger().SetLevel(core.LOG_DEBUG)

	}
	err=GctfDataManage.CreateTables(User{},Problems{},UserProblems{},Hints{},Tag{},Teams{})
	//GctfDataManage.DropTables(User{},Problems{},UserProblems{},Hints{},Tag{},Teams{})
	checkerr(err)
}

func checkerr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}
type User struct {
	Id             int64  `xorm:"autoincr pk 'id'"`
	Username       string `xorm:"unique"`
	Password       string
	Email          string    `xorm:"unique"`
	RegisterTime   time.Time `xorm:"created notnull"`
	ProblemsID     string
	SolvedProblems string
	Score          int
}
type Problems struct {
	Id          int64  `xorm:"autoincr pk 'id'"`
	Name        string `xorm:"unique"`
	Description string //Problem Description
	Value       int    //score
	Category    string
	Hidden      bool                   //should be problem hide?
	Location    string                 //saved physical position
	Scale       int `xorm:"default 0"` //score scale when each answer submit
}

type UserProblems struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	UserId     int64  //foreignkey gctf_user.id
	Location   string //problem net location
	Flag       string
	ProblemsId int64 //foreignkey gctf_problems.id
}

type Hints struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	ProblemsId int64 //foreignkey gctf_problems.id
	Hint       string
	Cost       int //cost score to get hint
}

type Tag struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	ProblemsId int64 //foreignkey gctf_problems.id
	Tag        string
}

type Teams struct {
	Id     int64  `xorm:"autoincr pk 'id'"`
	Name   string `xorm:"unique"`
	Member string
	Banned bool //if true this team can't login
}
