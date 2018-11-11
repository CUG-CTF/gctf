package model

import (
	"../gctfConfig"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

// meaning GctfDataMange is a database
var GctfDataManage *xorm.Engine

func init() {
	var err error
	// go-xorm is used to create database engine
	// engine, err := xorm.NewEngine(driverName, dataSourceName)
	// data
	GctfDataManage, err = xorm.NewEngine(gctfConfig.GCTF_DB_DRIVER, gctfConfig.GCTF_DB_STRING)
	// All table name have a gctf_ prefix

	// prefix，前缀
	tbMapper := core.NewPrefixMapper(core.GonicMapper{}, "gctf_")

	// fix problem_I_D to problem_id
	GctfDataManage.SetColumnMapper(core.GonicMapper{})
	GctfDataManage.SetTableMapper(tbMapper)

	// Ping is test the database is alive
	err = GctfDataManage.Ping()
	if err != nil {
		log.Fatal("database connect error:", err.Error())
	}
	if gctfConfig.GCTF_DEBUG {
		//GctfDataManage.ShowSQL(true)
		//GctfDataManage.Logger().SetLevel(core.LOG_DEBUG)

	}
	// this is create lots of tables?
	err = GctfDataManage.CreateTables(User{}, Problems{}, UserProblems{}, Hints{}, Tag{}, Teams{})
	// GctfDataManage.DropTables("gctf_user","gctf_problems","gctf_user_problems","gctf_hints","gctf_tag","gctf_teams")
	checkerr(err)
}

func init() {

}

func checkerr(err error) {
	if err != nil {
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
	IsAdmin        bool `xorm:"'is_admin'"`
}

type Problems struct {
	Id          int64  `xorm:"autoincr pk 'id'"`
	Name        string `xorm:"unique"`
	Description string // Problem Description
	Value       int    // score
	Category    string
	Hidden      bool                   // should be problem hide?
	Location    string                 // saved physical position
	Scale       int `xorm:"default 0"` // score scale when each answer submit
}

type UserProblems struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	UserId     int64  // foreignkey gctf_user.id
	Location   string // problem net location
	Flag       string
	ProblemsId int64 // foreignkey gctf_problems.id
}

type Hints struct {
	Id         int64 `xorm:"autoincr pk 'id'"`
	ProblemsId int64 // foreignkey gctf_problems.id
	Hint       string
	Cost       int // cost score to get hint
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
	Banned bool   // if true this team can't login
	Token  string // team token
}
